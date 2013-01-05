#!/usr/bin/env bash
URL=http://www.mdbg.net/chindict/export/cedict/cedict_1_0_ts_utf-8_mdbg.txt.gz
mkdir -p dict && cd dict && wget -c $URL
