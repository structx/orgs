
-- name: CreateOrganization :one
INSERT INTO organizations (
    name
) VALUES (
    $1
) RETURNING *;