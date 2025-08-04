-- name: InsertUser :execresult
INSERT INTO users (
    username,
    hashed_password 
) VALUES (
  ?, ?
);

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = ?;

-- name: InsertFile :execresult
INSERT INTO files (
  user_id,
  filename,
  content_type,
  size
) VALUES (
  ?, ?, ?, ?
);
