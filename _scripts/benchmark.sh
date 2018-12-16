#!/bin/sh
. ./.env
bombardier -m GET -l -H "${Authorization}" -n 100000 http://localhost:3000/api/v1/bookmarks
