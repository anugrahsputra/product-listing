# Product Listing API

A high-performance backend for product management, built with Go. This API provides comprehensive endpoints for managing categories and products with a focus on speed, reliability, and clean architecture.

## 🚀 Features

- **Categorized Products**: Organise products into distinct categories.
- **Unique Slug Generation**: SEO-friendly URL slugs for both categories and products.
- **CRUD Operations**: Complete Create, Read, Update, and Delete capabilities.
- **Optimized Performance**: Built on Gin and pgx/v5 for maximum throughput.
- **Type-Safe Database Access**: Powered by `sqlc` for compile-time verified SQL.
- **Clean Architecture**: Decoupled layers (Delivery, Usecase, Repository, Domain) for maintainability.

## 🛠 Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin (HTTP)
- **Database**: PostgreSQL
- **DB Driver**: pgx/v5
- **SQL Generator**: [sqlc](https://sqlc.dev/)
- **Configuration**: [cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- **Logging**: [go-logging](https://github.com/op/go-logging)

## 📋 Prerequisites

- [Go](https://go.dev/doc/install) 1.25 or higher
- [PostgreSQL](https://www.postgresql.org/download/) database
- [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) (for development)

## ⚙️ Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Update the `.env` file with your database credentials and preferred port.

## 🏃 Getting Started

### Installation

```bash
go mod download
```

### Build and Run

```bash
# Build the binary
go build -o bin/api ./cmd/api

# Run the API
./bin/api
```

Alternatively, run directly:
```bash
go run ./cmd/api
```

## 🔌 API Endpoints

### Categories
- `GET /api/category` - List all categories (with pagination)
- `GET /api/category/:id` - Get category by ID
- `GET /api/category/slug/:slug` - Get category by slug
- `POST /api/category` - Create a new category
- `PUT /api/category/:id` - Update an existing category
- `DELETE /api/category/:id` - Delete a category

### Products
- `GET /api/products/` - List all products (with pagination)
- `GET /api/products/:id` - Get product by ID
- `GET /api/products/category/:category_id` - List products in a specific category
- `POST /api/products/` - Create a new product
- `PUT /api/products/:id` - Update an existing product
- `DELETE /api/products/:id` - Delete a product

## 🧪 Development & Testing

### Seeding Data
Populate your database with 10 categories and 100 products:
```bash
./seed.sh
```

### Performance Testing
Run a stress test simulating 10,000 requests across mixed endpoints:
```bash
go run cmd/stress/main.go
```

### Verification
Run the full endpoint verification suite:
```bash
./verify.sh
```

## 📁 Project Structure

```
├── cmd/
│   ├── api/          # Main application entry point
│   └── stress/       # Load testing tool
├── internal/
│   ├── delivery/     # HTTP Handlers, DTOs, and Routing
│   ├── domain/       # Core Business Entities and Interfaces
│   ├── usecase/      # Business Logic implementation
│   ├── repository/   # Data Access implementation
│   └── db/           # Generated SQL code (sqlc)
├── sql/
│   ├── queries/      # SQL query definitions
│   └── schema/       # Database migrations
└── pkg/
    └── logger/       # Shared logging utilities
```

## 📜 License

This project is licensed under the MIT License.
