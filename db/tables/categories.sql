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

INSERT INTO categories (name, slug) VALUES ('Flu-like', 'flu-like');
INSERT INTO categories (name, slug) VALUES ('Gastrointestinal', 'gastrointestinal');
INSERT INTO categories (name, slug) VALUES ('Psychological', 'psychological');
INSERT INTO categories (name, slug) VALUES ('Life threatening', 'life-threatening');
INSERT INTO categories (name, slug) VALUES ('Skin & localized to injection site', 'skin-and-localized-to-injection-site');
INSERT INTO categories (name, slug) VALUES ('Muscles & bones', 'muscles-and-bones');
INSERT INTO categories (name, slug) VALUES ('Immune system & inflammation', 'immune-system-and-inflammation');
INSERT INTO categories (name, slug) VALUES ('Nervous system', 'nervous-system');
INSERT INTO categories (name, slug) VALUES ('Cardiovascular', 'cardiovascular');
INSERT INTO categories (name, slug) VALUES ('Eyes, mouth & ears', 'eyes-mouth-and-ears');
INSERT INTO categories (name, slug) VALUES ('Errors by medical staff', 'errors-by-medical-staff');
INSERT INTO categories (name, slug) VALUES ('Breathing', 'breathing');
INSERT INTO categories (name, slug) VALUES ('Urinary', 'urinary');
INSERT INTO categories (name, slug) VALUES ('Balance & mobility', 'balance-and-mobility');
INSERT INTO categories (name, slug) VALUES ('Gynecological', 'gynecological');
