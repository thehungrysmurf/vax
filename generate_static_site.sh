#!/bin/bash
set -e

mkdir -p docs/vaccine/pfizer

cp -r assets docs

curl http://localhost:8888 > docs/index.html

curl http://localhost:8888/vaccine/pfizer > docs/vaccine/pfizer/index.html
