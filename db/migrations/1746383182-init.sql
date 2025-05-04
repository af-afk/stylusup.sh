-- migrate:up

CREATE TABLE stylusup_popularity_1 (
	id SERIAL PRIMARY KEY,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ip_hash CHAR(64) NOT NULL,
	lang CHAR(100) NOT NULL,
	arch CHAR(32) NOT NULL,
	os CHAR(100) NOT NULL
);

CREATE INDEX ON stylusup_popularity_1 (lang, arch, os);

CREATE UNIQUE INDEX ON stylusup_popularity_1 (ip_hash);

-- migrate:down
