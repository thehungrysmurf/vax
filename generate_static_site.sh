#!/bin/bash
set -e

mkdir -p docs/about
cp -r assets docs

VACCINES=(pfizer moderna janssen)
CATEGORIES=(flu-like gastrointestinal psychological life-threatening skin-and-localized-to-injection-site muscles-and-bones immune-system-and-inflammation nervous-system cardiovascular eyes-mouth-and-ears urinary breathing balance-and-mobility gynecological)
SEXES=(male female)
AGE_GROUPS=(12/15 16/25 26/39 40/59 60/75 76/89 90/110)
VAX_HOST=http://localhost:8888

for VACCINE in ${VACCINES[@]}; do
  VACCINE_PATH=vaccine/$VACCINE

  mkdir -p docs/$VACCINE_PATH
  echo ">> $VAX_HOST/$VACCINE_PATH/ > docs/$VACCINE_PATH/index.html"
  curl -s $VAX_HOST/$VACCINE_PATH/ > docs/$VACCINE_PATH/index.html

  for SEX in ${SEXES[@]}; do
    for CATEGORY in ${CATEGORIES[@]}; do
      if [[ "$CATEGORY" = "gynecological" && "$SEX" = "male" ]]; then
        continue
      fi
      for AGE_GROUP in ${AGE_GROUPS[@]}; do
        SUMMARY_PATH=vaccine/$VACCINE/category/$CATEGORY/$SEX/$AGE_GROUP
        mkdir -p docs/$SUMMARY_PATH
        echo ">> $VAX_HOST/$SUMMARY_PATH/ > docs/$SUMMARY_PATH/index.html"
        curl -s $VAX_HOST/$SUMMARY_PATH/ > docs/$SUMMARY_PATH/index.html
      done
    done
  done
done

curl -s $VAX_HOST/ > docs/index.html
curl -s $VAX_HOST/about/ > docs/about/index.html
curl -s $VAX_HOST/404 > docs/404.html

