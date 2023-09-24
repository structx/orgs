
-- name: CreateOrganization :one
-- CreateOrganization creates a new organization
INSERT INTO organizations (
    name,
    processor_id
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetOrganization :one
-- GetOrganization retrieves an organization by id
SELECT * FROM organizations WHERE id = $1;

-- name: ListOrganizations :many
-- ListOrganizations retrieves all organizations
SELECT * FROM organizations;

-- name: UpdateOrganization :one
-- UpdateOrganization updates an organization by id
UPDATE organizations SET name = $1 WHERE id = $2 RETURNING *;

-- name: DeleteOrganization :exec
-- DeleteOrganization deletes an organization by id
DELETE FROM organizations WHERE id = $1;