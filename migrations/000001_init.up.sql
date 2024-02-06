--
-- Enable extensions
--
CREATE EXTENSION IF NOT EXISTS vector;


--
-- Actors
--
CREATE TABLE actors(
	id SERIAL PRIMARY KEY,
	name VARCHAR(64) NOT NULL,
	embedding vector(768)
);


--
-- Conversations
--
CREATE TABLE conversations(
	id SERIAL PRIMARY KEY,
	title VARCHAR(64) NOT NULL,
	description VARCHAR(64) NOT NULL,
	embedding vector(768)
);


--
-- Conversation lines
--
CREATE TABLE lines (
	id SERIAL PRIMARY KEY,
	conversation_id INT REFERENCES conversations(id) NOT NULL,
	actor_id INT REFERENCES actors(id) NOT NULL,
	stamp TIMESTAMP NOT NULL,
	content TEXT NOT NULL,
	embedding vector(768)
);
