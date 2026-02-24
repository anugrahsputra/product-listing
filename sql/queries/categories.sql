-- name: CreateCategory :one
INSERT INTO categories (name, slug, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id, name, slug, created_at, updated_at;

-- name: GetCategories :many
SELECT id, name, slug, created_at, updated_at
FROM categories
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: GetCategoryById :one
SELECT id, name, slug, created_at, updated_at
FROM categories
WHERE id = $1;

-- name: GetCategoryBySlug :one
SELECT id, name, slug, created_at, updated_at
FROM categories
WHERE slug = $1;

-- name: UpdateCategory :exec
UPDATE categories
SET
    name = COALESCE($2, name),
    slug = COALESCE($3, slug),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: GetCategoriesCount :one
SELECT COUNT(*) FROM categories;

