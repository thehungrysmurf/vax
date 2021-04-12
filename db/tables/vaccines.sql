DROP TABLE IF EXISTS vaccines CASCADE;

DROP TYPE IF EXISTS ILLNESS;
CREATE TYPE ILLNESS AS ENUM ('covid19');

CREATE TABLE vaccines(

	id SERIAL
		PRIMARY KEY,

	illness ILLNESS
		NOT NULL,

	manufacturer VARCHAR(255)
		NOT NULL,

	created_at TIMESTAMPTZ
		NOT NULL
        DEFAULT NOW()
);

INSERT INTO vaccines (illness, manufacturer) VALUES ('covid19', 'moderna');
INSERT INTO vaccines (illness, manufacturer) VALUES ('covid19', 'pfizer');
INSERT INTO vaccines (illness, manufacturer) VALUES ('covid19', 'janssen');
