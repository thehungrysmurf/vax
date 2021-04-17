DROP TABLE IF EXISTS vaccination_totals;

CREATE TABLE vaccination_totals(

	id BIGSERIAL
	    PRIMARY KEY,

	pfizer BIGINT
		NOT NULL,

	moderna BIGINT
		NOT NULL,

	janssen BIGINT
	    NOT NULL,

	created_at TIMESTAMPTZ
		NOT NULL
        DEFAULT NOW(),

    updated_at TIMESTAMPTZ
        NOT NULL
);