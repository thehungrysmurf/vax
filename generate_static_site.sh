#!/bin/bash
set -e

mkdir -p static_site/vaccine/pfizer

cp -r assets static_site

curl http://localhost:8888 > static_site/index.html

curl http://localhost:8888/vaccine/pfizer > static_site/vaccine/pfizer/index.html
