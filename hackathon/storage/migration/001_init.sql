-- +migrate Up
CREATE TABLE users (
  id bigint AUTO_INCREMENT PRIMARY KEY,
  username nvarchar(255) NOT NULL,
  hashed_password nvarchar(255) NOT NULL
);

CREATE UNIQUE INDEX idx_users_username
ON users (username);

-- +migrate Down
DROP TABLE IF EXISTS users;
