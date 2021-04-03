DROP TABLE IF EXISTS people_symptoms;

CREATE TABLE people_symptoms(

	vaers_id BIGSERIAL
		NOT NULL,

	symptom_id BIGSERIAL
		NOT NULL,

	vaccine_id SERIAL
		NOT NULL,

	FOREIGN KEY (vaers_id) REFERENCES people(vaers_id),
	FOREIGN KEY (symptom_id) REFERENCES symptoms(id),
	FOREIGN KEY (vaccine_id) REFERENCES vaccines(id),

    PRIMARY KEY (vaers_id, symptom_id, vaccine_id)
);
