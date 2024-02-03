
--
-- Enable extensions
--
CREATE EXTENSION vector;


--
-- Users
--
CREATE TABLE users(
	id SERIAL PRIMARY KEY,
	username VARCHAR(64) NOT NULL,
	embedding vector(1024) NOT NULL
);


--
-- Conversations
--
CREATE TABLE conversations(
	id SERIAL PRIMARY KEY,
	name VARCHAR(64) NOT NULL,
	embedding vector(1024) NOT NULL
);


--
-- Conversation lines
--
CREATE TABLE lines (
	id SERIAL PRIMARY KEY,
	conversation_id INT REFERENCES conversations(id) NOT NULL,
	user_id INT REFERENCES users(id) NOT NULL,
	stamp TIMESTAMP NOT NULL,
	content TEXT NOT NULL,
	embedding vector(1024) NOT NULL
);
