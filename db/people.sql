-- noinspection SqlNoDataSourceInspectionForFile

DROP TABLE IF EXISTS people;

CREATE TYPE SEX AS ENUM ('M', 'F', 'unknown');

CREATE TABLE people(

	id BIGSERIAL
		PRIMARY KEY,

	vaers_id BIGINT
		NOT NULL,

	age INT
		NOT NULL
		DEFAULT 0,

	sex SEX
		NOT NULL
		DEFAULT 'unknown',

	notes VARCHAR(512)
		NOT NULL
		DEFAULT '',

	reported_at TIMESTAMPTZ
		NOT NULL
        DEFAULT NOW(),

    created_at TIMESTAMPTZ
    	NOT NULL
        DEFAULT NOW()
);
