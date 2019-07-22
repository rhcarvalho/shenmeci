shenmeci
========

Shenmeci is a simple Chinese-English online dictionary.

It uses a dictionary-based Chinese word segmentation approach, matching maximal-length words from left to right.

Full-Text Search for English terms is also available.


Online Dictionary
-----------------

An online instance can be found at http://shenmeci.rodolfocarvalho.net.


Development repository
----------------------

The project source-code is hosted at https://github.com/rhcarvalho/shenmeci.


Development/deployment dependencies
-----------------------------------

* Go
* SQLite (>= [3.9.0][sqlite-390], for the [`json1`][sqlite-json1] extension)
* go-sqlite3
* DAWGo

[sqlite-390]: https://www.sqlite.org/releaselog/3_9_0.html
[sqlite-json1]: https://www.sqlite.org/json1.html

Running
-------

Download, compile and install in one go:

    $ go get -u -tags 'sqlite_json1 sqlite_fts5' github.com/rhcarvalho/shenmeci

Download the CEDICT Chinese-English dictionary:

    $ ./download_dict.sh

Before running Shenmeci you will need to create a configuration file like this:

    {
      "Http": {
        "Host": "127.0.0.1",
        "Port": 8080
      },
      "StaticPath": "static/",
      "CedictPath": "dict/cedict_1_0_ts_utf-8_mdbg.txt.gz"
    }

Make sure you have a new version of the SQLite library.
Start Shenmeci:

    $ shenmeci -config path/to/config.json


Migrating query logs from MongoDB to SQLite
-------------------------------------------

Data stored in MongoDB can be migrated using `mongoexport`, and imported using a
script:

    $ mongoexport --db shenmeci --collection queries --out shenmeci-export.json
    $ python3 import-queries.py --from shenmeci-export.json --to PATH/TO/shenmeci.sqlite
