#!/bin/bash

BASE_URL="http://localhost:8080/api"
PREFIX=$(date +%s)

echo "Creating 10 categories with prefix $PREFIX..."
for i in {1..10}
do
  curl -s -X POST "$BASE_URL/category" -H "Content-Type: application/json" -d "{\"name\": \"MCat $PREFIX $i\", \"slug\": \"mcat-$PREFIX-$i\"}" > /dev/null
done
sleep 1

echo "Fetching category IDs..."
# Increased limit to 1000
CATEGORY_IDS=$(curl -s "$BASE_URL/category?limit=1000" | jq -r ".data[] | select(.slug | contains(\"$PREFIX\")) | .id")
IDS_ARRAY=($CATEGORY_IDS)
NUM_IDS=${#IDS_ARRAY[@]}
echo "Found $NUM_IDS new categories"

if [ $NUM_IDS -eq 0 ]; then
    echo "No categories found. Check server output."
    exit 1
fi

echo "Creating 10 products, each in 2 categories..."
for i in {1..10}
do
  IDX1=$(( (i-1) % NUM_IDS ))
  IDX2=$(( (i) % NUM_IDS ))
  CAT_ID1=${IDS_ARRAY[$IDX1]}
  CAT_ID2=${IDS_ARRAY[$IDX2]}

  curl -s -X POST "$BASE_URL/products/" -H "Content-Type: application/json" -d "{\"name\": \"Multi-Cat Product $PREFIX $i\", \"slug\": \"prod-multi-$PREFIX-$i\", \"Description\": \"desc\", \"category_ids\": [\"$CAT_ID1\", \"$CAT_ID2\"], \"price\": 49.99}" > /dev/null
done

echo "--------------------------------"
echo "Verification:"
curl -s "http://localhost:8080/api/products/?limit=1" | jq '.data[0] | {name, categories}'
