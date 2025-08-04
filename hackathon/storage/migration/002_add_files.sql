-- +migrate Up
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

-- +migrate Down
DROP TABLE IF EXISTS files;
