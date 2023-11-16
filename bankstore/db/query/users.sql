-- name: CreateUser :one
INSERT INTO users (
    username,
    password,
    repassword   
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserByName :one
SELECT * FROM users 
WHERE username=$1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE username=$1;