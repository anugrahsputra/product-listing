# Learning Notes: Product Listing Backend

This document tracks the architectural decisions, lessons learned, and future concepts for the Product Listing API.

## ✅ Phase 1: The Foundation (One-to-Many Relationship)
*Status: Completed & Verified*

### What I've Learned
1.  **Clean Architecture**: How to separate code into Domain (entities), Usecase (logic), Repository (data), and Delivery (HTTP/JSON).
2.  **SQLC & Type Safety**: Using `sqlc` to generate Go code from raw SQL. It prevents runtime errors by verifying queries at compile time.
3.  **One-to-Many (1:N)**: Storing a `category_id` directly in the `products` table. This is perfect when a product belongs to exactly one category.
4.  **The "N+1" Problem & Joins**: 
    *   *Problem*: Fetching a list of products and then making a separate request for each image (slow).
    *   *Solution*: Using `LEFT JOIN` in a single SQL query to "attach" the primary image to the product data immediately.
5.  **Data Integrity**: Using Foreign Keys (`REFERENCES`) and `ON DELETE CASCADE` or `RESTRICT` to ensure data stays consistent.
6.  **Idempotency**: Adding `IF NOT EXISTS` to SQL scripts so they can be run multiple times without crashing the server.

---

## ✅ Phase 2: The Next Level (Many-to-Many Relationship)
*Status: Completed & Verified*

### What I've Learned
1.  **Join Tables (Junction Tables)**: 
    *   How to create a third table (e.g., `product_categories`) to bridge two entities.
    *   Removing the single-seat `category_id` from the main table to allow infinite connections.
2.  **PostgreSQL JSON Aggregation**:
    *   `jsonb_build_object()`: Turning database rows into structured JSON objects inside the query.
    *   `json_agg()`: Collapsing multiple related rows into a single JSON array field.
3.  **Go JSON Unmarshaling**:
    *   Using `json.Unmarshal` within the Repository layer to convert database JSON strings back into Go slices of structs.
4.  **Business Logic for Sets**:
    *   Clearing old relationships and adding new ones during an "Update" operation to maintain the many-to-many set.

---

## 🚀 Phase 3: Future Concepts
*Status: Planned*

### What I Need to Learn
1.  **Transaction Management**: Ensuring that when we create a product AND its category links, either both succeed or both fail (Atomic operations).
2.  **Indexing JSONB**: How to optimize queries that filter or search inside JSON columns.
3.  **Custom Marshalling**: Implementing the `Unmarshaler` interface in Go to automate JSON conversion.

