-- name: CreateProduct :one
INSERT INTO products(name, slug, description, price, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, name, slug, description, price, created_at, updated_at;

-- name: GetAllProducts :many
SELECT 
    p.id,
    p.name,
    p.slug,
    p.description,
    p.price,
    p.created_at,
    p.updated_at,
    pi.url as primary_image_url,
    (
        SELECT json_agg(jsonb_build_object('id', c.id, 'name', c.name, 'slug', c.slug))
        FROM product_categories pc
        JOIN categories c ON c.id = pc.category_id
        WHERE pc.product_id = p.id
    )::json as categories
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
    p.price,
    p.created_at,
    p.updated_at,
    pi.url as primary_image_url,
    (
        SELECT json_agg(jsonb_build_object('id', c.id, 'name', c.name, 'slug', c.slug))
        FROM product_categories pc
        JOIN categories c ON c.id = pc.category_id
        WHERE pc.product_id = p.id
    )::json as categories
FROM products p
LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
WHERE p.id = $1;

-- name: GetProductsByCategoryID :many
SELECT 
    p.id,
    p.name,
    p.slug,
    p.description,
    p.price,
    p.created_at,
    p.updated_at,
    pi.url as primary_image_url,
    (
        SELECT json_agg(jsonb_build_object('id', c.id, 'name', c.name, 'slug', c.slug))
        FROM product_categories pc_all
        JOIN categories c ON c.id = pc_all.category_id
        WHERE pc_all.product_id = p.id
    )::json as categories
FROM products p
JOIN product_categories pc ON p.id = pc.product_id
LEFT JOIN product_images pi ON p.id = pi.product_id AND pi.is_primary = true
WHERE pc.category_id = $1;

-- name: AddProductCategory :exec
INSERT INTO product_categories (product_id, category_id)
VALUES ($1, $2);

-- name: ClearProductCategories :exec
DELETE FROM product_categories
WHERE product_id = $1;

-- name: UpdateProduct :exec
UPDATE products
SET 
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    price = COALESCE($4, price),
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: GetProductsCount :one
SELECT COUNT(*) FROM products;
