#!/bin/bash

# ping write-file read-file read-write-file
k6 run -e MY_HOSTNAME=$1 -d $2 -u $3 ./test.js
