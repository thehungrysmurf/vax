DROP TABLE IF EXISTS categories CASCADE;

CREATE TABLE categories(

	id SERIAL
		PRIMARY KEY,

	name VARCHAR(255)
		NOT NULL,

	slug VARCHAR(255)
	    NOT NULL,

	created_at TIMESTAMPTZ
		NOT NULL
        DEFAULT NOW()
);

INSERT INTO categories (name, slug) VALUES ('flu-like', 'flu-like');
INSERT INTO categories (name, slug) VALUES ('gastrointestinal', 'gastrointestinal');
INSERT INTO categories (name, slug) VALUES ('psychological', 'psychological');
INSERT INTO categories (name, slug) VALUES ('life threatening', 'life-threatening');
INSERT INTO categories (name, slug) VALUES ('skin & localized to injection site', 'skin-and-localized-to-injection-site');
INSERT INTO categories (name, slug) VALUES ('muscles & bones', 'muscles-and-bones');
INSERT INTO categories (name, slug) VALUES ('allergic', 'allergic');
INSERT INTO categories (name, slug) VALUES ('nervous system', 'nervous-system');
INSERT INTO categories (name, slug) VALUES ('cardiovascular', 'cardiovascular');
INSERT INTO categories (name, slug) VALUES ('eyes, mouth & ears', 'eyes-mouth-and-ears');
