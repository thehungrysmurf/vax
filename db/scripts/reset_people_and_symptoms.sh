#!/bin/bash

psql -d vax -f ./drop_people_and_symptoms.sql

psql -d vax -f ../tables/people.sql
psql -d vax -f ../tables/symptoms.sql
psql -d vax -f ../tables/people_symptoms.sql
psql -d vax -f ../tables/symptoms_categories.sql