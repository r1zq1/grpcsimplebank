-- name: CreateTransfer :one
INSERT INTO transfers(from_account_id, to_account_id, amount) 
VALUES ($1, $2, $3) 
RETURNING *;

-- name: ListTransfersByAccountID :many
SELECT * FROM transfers
WHERE from_account_id = $1 OR to_account_id = $1
ORDER BY id DESC
LIMIT 100;