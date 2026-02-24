-- name: CreateProductImage :one
INSERT INTO product_images (
    product_id,
    url,
    is_primary
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetProductImages :many
SELECT * FROM product_images
WHERE product_id = $1
ORDER BY is_primary DESC, created_at ASC;

-- name: GetProductPrimaryImage :one
SELECT * FROM product_images
WHERE product_id = $1 AND is_primary = true
LIMIT 1;

-- name: DeleteProductImage :exec
DELETE FROM product_images
WHERE id = $1;

-- name: SetProductPrimaryImage :exec
UPDATE product_images
SET is_primary = (id = $2)
WHERE product_id = $1;
