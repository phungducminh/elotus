create database elotus;
use elotus;

CREATE TABLE users (
  id bigint AUTO_INCREMENT PRIMARY KEY,
  username nvarchar(255) NOT NULL,
  hashed_password nvarchar(255) NOT NULL
);

CREATE UNIQUE INDEX idx_users_username
ON users (username);

CREATE TABLE files (
  id bigint AUTO_INCREMENT PRIMARY KEY,
  user_id bigint NOT NULL,
  filename nvarchar(255) NOT NULL,
  content_type nvarchar(255) NOT NULL,
  size integer NOT NULL,
  FOREIGN KEY (user_id)
      REFERENCES users(id)
      ON DELETE CASCADE
);
