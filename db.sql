CREATE TABLE records
(
	id SERIAL PRIMARY KEY NOT NULL,
	user_id INT NOT NULL,
	body VARCHAR(128) NOT NULL,
	date_added TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);

CREATE INDEX idx_user_id ON records (user_id);
CREATE INDEX idx_date_added ON records (date_added);
CREATE INDEX idx_body ON records (body);
