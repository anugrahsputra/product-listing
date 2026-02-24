-- name: CreateProduct :one
INSERT INTO products(name, slug, description, category_id, price, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id, name, slug, description, category_id, price, created_at, updated_at;

-- name: GetAllProducts :many
SELECT 
    p.id,
    p.name,
    p.slug,
    p.description,
    p.category_id,
    p.price,
    p.created_at,
    p.updated_at,
    pi.url as primary_image_url
FROM products p
LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
ORDER BY p.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetProductByID :one
SELECT 
    p.id,
    p.name,
    p.slug,
    p.description,
    p.category_id,
    p.price,
    p.created_at,
    p.updated_at,
    pi.url as primary_image_url
FROM products p
LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
WHERE p.id = $1;

-- name: GetProductByCategory :many
SELECT 
    p.id,
    p.name,
    p.slug,
    p.description,
    p.category_id,
    p.price,
    p.created_at,
    p.updated_at,
    pi.url as primary_image_url
FROM products p
LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
WHERE p.category_id = $1;


-- name: UpdateProduct :exec
UPDATE products
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    category_id = COALESCE($4, category_id),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: GetProductsCount :one
SELECT COUNT(*) FROM products;

