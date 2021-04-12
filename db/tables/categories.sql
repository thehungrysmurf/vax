DROP TABLE IF EXISTS categories CASCADE;

CREATE TABLE categories(

	id SERIAL
		PRIMARY KEY,

	name VARCHAR(255)
		NOT NULL,

	created_at TIMESTAMPTZ
		NOT NULL
        DEFAULT NOW()
);

INSERT INTO categories (name) VALUES ('flu-like');
INSERT INTO categories (name) VALUES ('gastrointestinal');
INSERT INTO categories (name) VALUES ('psychological');
INSERT INTO categories (name) VALUES ('life threatening');
INSERT INTO categories (name) VALUES ('skin & localized to injection site');
INSERT INTO categories (name) VALUES ('muscles & bones');
INSERT INTO categories (name) VALUES ('allergic');
INSERT INTO categories (name) VALUES ('nervous system');
INSERT INTO categories (name) VALUES ('cardiovascular');
INSERT INTO categories (name) VALUES ('eyes, mouth & ears');
