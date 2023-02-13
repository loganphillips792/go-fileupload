-- DROP TABLE IF EXISTS images;

-- CREATE TABLE IF NOT EXISTS images (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, file_path TEXT);
-- CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, email TEXT, password TEXT);
-- CREATE TABLE IF NOT EXISTS sessions (session_id TEXT PRIMARY KEY, expires_at INTEGER);

-- INSERT INTO images (name, file_path) VALUES ('Name 1', 'Name 1');
-- INSERT INTO images (name, file_path) VALUES ('Name 2', 'Name 2');
-- INSERT INTO images (name, file_path) VALUES ('Name 3', 'Name 3');

-- DROP TABLE IF EXISTS images;
-- DROP TABLE IF EXISTS users;
-- DROP TABLE IF EXISTS sessions;

-- Create 'images' table
CREATE TABLE IF NOT EXISTS images (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  file_path VARCHAR(255) NOT NULL
);

-- Create 'users' table
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(50) NOT NULL,
  email VARCHAR(255) NULL,
  password VARCHAR(255) NOT NULL
);

-- Create 'sessions' table
CREATE TABLE IF NOT EXISTS sessions (
  session_id TEXT PRIMARY KEY,
  expires_at TIMESTAMP NOT NULL
);

-- Insert test data into 'images' table
INSERT INTO images (name, file_path)
VALUES ('image1', '/path/to/image1.png'),
       ('image2', '/path/to/image2.png');

-- Insert test data into 'users' table
INSERT INTO users (username, email, password)
VALUES ('user1', 'user1@example.com', 'password1'),
       ('user2', 'user2@example.com', 'password2');

-- Insert test data into 'sessions' table
INSERT INTO sessions (session_id, expires_at)
VALUES ('6fcbcf6c-7366-4ec6-9b84-11f13d6308f9', '2023-02-13 23:59:59'),
       ('edf7b15d-6d0b-4272-a1a3-3049bde67d15', '2023-02-14 23:59:59');
