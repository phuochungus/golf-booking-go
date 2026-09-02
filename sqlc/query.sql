-- name: CreateFacility :exec
INSERT INTO facilities (name, organization_id) VALUES (?, ?);
-- name: GetFacilities :many
SELECT * 
FROM facilities
WHERE 
    1 = 1
    AND (sqlc.narg('lte_rating') IS NULL OR rating <= sqlc.narg('lte_rating'))
    AND (sqlc.narg('gte_rating') IS NULL OR rating >= sqlc.narg('gte_rating'))
    AND (sqlc.narg('organization_id') IS NULL OR organization_id = sqlc.narg('organization_id'));

-- name: GetFacilityByID :one
SELECT * FROM facilities WHERE id = ?;
-- name: UpdateFacility :exec
UPDATE facilities 
SET 
    name = coalesce(sqlc.narg('name'), name), 
    organization_id = coalesce(sqlc.narg('organization_id'), organization_id) 
WHERE id = sqlc.arg('id');
