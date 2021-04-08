DROP TABLE IF EXISTS symptoms_categories;

CREATE TABLE symptoms_categories(

	symptom_id BIGINT
		NOT NULL,

	category_id BIGINT
		NOT NULL,

	FOREIGN KEY (symptom_id) REFERENCES symptoms(id),
    FOREIGN KEY (category_id) REFERENCES categories(id),

	PRIMARY KEY (symptom_id, category_id)
);
