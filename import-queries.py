import sqlite3
import argparse
import sys


parser = argparse.ArgumentParser(
    description="Import Shenmeci Query data from a MongoDB export to SQLite."
)
parser.add_argument("--from", required=True, help="path to MongoDB export")
parser.add_argument("--to", required=True, help="path to SQLite database")


def printerr(*value):
    print(*value, sep="\n", file=sys.stderr)


if __name__ == "__main__":
    args = vars(parser.parse_args())
    db = sqlite3.connect(args["to"])

    if db.execute("SELECT 1 FROM sqlite_master WHERE type='table' AND name='query'").fetchone():
        printerr(
            f"""Table 'query' already in '{args["to"]}'""",
            "Not safe to import: could add duplicates!",
            "Aborting...",
        )
        sys.exit(1)

    db.execute("CREATE TABLE query(json)")
    with db:
        # Open file, but do not decode bytes just yet -- some lines may not be
        # valid UTF-8.
        with open(args["from"], "rb") as f:
            def gen():
                for line in f:
                    try:
                        yield (line.decode("UTF-8"),)
                    except UnicodeDecodeError:
                        # Found some queries from Korea!
                        yield (line.decode("EUC-KR"),)
            c = db.executemany("INSERT INTO query VALUES(json(?))", gen())
            print(f"imported {c.rowcount} rows")

    print("ok")
