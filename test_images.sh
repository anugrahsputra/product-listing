#!/bin/bash

BASE_URL="http://localhost:8080/api"
TS=$(date +%s)

echo "--- 1. Setup: Creating temporary category and product ---"
curl -s -X POST "$BASE_URL/category" -H "Content-Type: application/json" -d "{\"name\": \"Img Test Cat $TS\", \"slug\": \"img-cat-$TS\"}" > /dev/null
CAT_ID=$(curl -s "$BASE_URL/category/slug/img-cat-$TS" | jq -r '.data.id')

curl -s -X POST "$BASE_URL/products/" -H "Content-Type: application/json" -d "{\"name\": \"Img Test Prod $TS\", \"slug\": \"img-prod-$TS\", \"Description\": \"test\", \"category_id\": \"$CAT_ID\", \"price\": 1.0}" > /dev/null
PROD_ID=$(curl -s "$BASE_URL/products/?limit=100" | jq -r ".data[] | select(.slug==\"img-prod-$TS\") | .id")

echo "Using Product ID: $PROD_ID"

echo "--- 2. Testing PRODUCT IMAGES Endpoints ---"

# Add Image 1
echo "POST /product-images (Image 1)..."
IMG1_JSON=$(curl -s -X POST "$BASE_URL/product-images" -H "Content-Type: application/json" -d "{\"product_id\": \"$PROD_ID\", \"url\": \"https://example.com/img1.jpg\", \"is_primary\": true}")
echo $IMG1_JSON | jq -c
IMG1_ID=$(echo $IMG1_JSON | jq -r '.data.id')

# Add Image 2
echo "POST /product-images (Image 2)..."
IMG2_JSON=$(curl -s -X POST "$BASE_URL/product-images" -H "Content-Type: application/json" -d "{\"product_id\": \"$PROD_ID\", \"url\": \"https://example.com/img2.jpg\", \"is_primary\": false}")
echo $IMG2_JSON | jq -c
IMG2_ID=$(echo $IMG2_JSON | jq -r '.data.id')

# Get Images
echo "GET /product-images/product/:product_id:"
curl -s "$BASE_URL/product-images/product/$PROD_ID" | jq -c '{status, count: (.data | length)}'

# Set Primary
echo "PUT /product-images/primary/:product_id/:image_id (Set Image 2 as primary):"
curl -s -X PUT "$BASE_URL/product-images/primary/$PROD_ID/$IMG2_ID" | jq -c

# Delete Image
echo "DELETE /product-images/:id (Deleting Image 1):"
curl -s -X DELETE "$BASE_URL/product-images/$IMG1_ID" | jq -c

echo "--- 3. Cleanup ---"
curl -s -X DELETE "$BASE_URL/products/$PROD_ID" > /dev/null
curl -s -X DELETE "$BASE_URL/category/$CAT_ID" > /dev/null

echo "--- ALL IMAGE TESTS COMPLETED ---"
