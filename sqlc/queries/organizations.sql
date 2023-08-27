
-- name: CreateOrganization :one
INSERT INTO organizations (
    name,
    processor_id
) VALUES (
    $1, $2
) RETURNING *;