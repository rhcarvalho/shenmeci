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


Running
-------

You will need a recent version of [Go](https://go.dev).

Checkout the repository and then:

1. Download the CEDICT Chinese-English dictionary:

    $ ./download_dict.sh

2. Create a configuration file in the repository root called `config.json`:

    {
      "Http": {
        "Host": "127.0.0.1",
        "Port": 8080
      },
      "CedictPath": "dict/cedict_1_0_ts_utf-8_mdbg.txt.gz"
    }

3. Start the HTTP server:

    $ go run -tags sqlite_fts5 .

## Deployment

```
bin/deploy-fly
```
