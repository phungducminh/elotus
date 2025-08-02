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
