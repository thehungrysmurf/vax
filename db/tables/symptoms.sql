DROP TABLE IF EXISTS symptoms CASCADE;

CREATE TABLE symptoms(

	id BIGSERIAL
		PRIMARY KEY,

	name VARCHAR(255)
		NOT NULL,

	alias VARCHAR(255)
		NOT NULL
		DEFAULT '',

	created_at TIMESTAMPTZ
		NOT NULL
        DEFAULT NOW(),

    UNIQUE(name)
);
