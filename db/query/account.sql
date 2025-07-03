-- name: GetAccountForUpdate :one
SELECT * FROM accounts 
WHERE id = $1 
FOR UPDATE;

-- name: AddAccountBalance :one
UPDATE accounts 
SET balance = balance + $2 
WHERE id = $1 
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts 
WHERE id = $1 
FOR UPDATE;
