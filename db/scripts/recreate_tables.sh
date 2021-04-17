#!/bin/bash

psql -d vax -f ./db/tables/people.sql
psql -d vax -f ./db/tables/symptoms.sql
psql -d vax -f ./db/tables/people_symptoms.sql
psql -d vax -f ./db/tables/symptoms_categories.sql
psql -d vax -f ./db/tables/vaccination_totals.sql