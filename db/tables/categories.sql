-- noinspection SqlNoDataSourceInspectionForFile

DROP TABLE IF EXISTS categories;

CREATE TABLE categories(

	id SERIAL
		PRIMARY KEY,

	category VARCHAR(255)
		NOT NULL,

	created_at TIMESTAMPTZ
		NOT NULL
        DEFAULT NOW()
);
