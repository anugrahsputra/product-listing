# AGENTS.md - Development Guide for AI Agents

This document provides guidelines for working on the product-listing backend codebase.

## Philosophy & Guidelines

### Core Philosophy

- **Safety First**
  Never risk user data, stability, or backward compatibility.
  When uncertain, stop and ask for clarification.

- **Incremental Progress**
  Break complex tasks into small, verifiable steps.
  Large, speculative changes are forbidden.

- **Clear Intent Over Cleverness**
  Prefer readable, boring, maintainable solutions.
  Clever hacks are a liability.

- **Native Performance Mindset**
  Optimize only when necessary and with evidence.
  Avoid premature optimization.

---

### Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:

- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

### Simplicity first

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

### Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:

- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

### Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:

- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:

```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

---

## Project Overview

- **Language**: Go 1.25+
- **Framework**: Gin (HTTP), pgx/v5 (PostgreSQL)
- **Architecture**: Clean Architecture (domain, usecase, repository, delivery layers)
- **Database**: PostgreSQL with sqlc for type-safe SQL

## Build, Lint, and Test Commands

### Build

```bash
go build -o bin/api ./cmd/api
```

### Run

```bash
go run ./cmd/api
# or with custom env
PORT=8081 DATABASE_URL="postgres://user:pass@localhost:5432/db" go run ./cmd/api
```

### Test

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./internal/domain/...

# Run a single test function
go test -v -run TestFunctionName ./internal/domain/...

# Run tests with coverage
go test -cover ./...
```

### Code Generation

```bash
# Regenerate SQL code (after modifying sql/queries/*.sql or sql/schema/*.sql)
sqlc generate
```

### Dependencies

```bash
go mod tidy
go mod download
```

## Code Style Guidelines

### Project Structure

```
cmd/api/main.go           # Application entry point
config/                   # Configuration
internal/
  domain/                 # Entities and repository interfaces
  usecase/                # Business logic (interfaces + implementations)
  repository/             # Data access layer implementations
  delivery/
    handler/              # HTTP handlers
    dto/                  # Data Transfer Objects
    router/               # Route definitions
pkg/logger/               # Shared logger package
sql/
  queries/                # SQL query files
  schema/                 # Database schema
```

### Imports

- Standard library imports first
- Third-party imports second (github.com, etc.)
- Project imports last
- Group imports with blank line between groups
- Use canonical import paths (e.g., `product-listing/internal/domain`)

Example:

```go
import (
    "context"
    "fmt"
    "net/http"
    "product-listing/internal/domain"
    "product-listing/internal/delivery/dto"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)
```

### Naming Conventions

**Files**: Use snake_case (e.g., `product_handler.go`, `product_usecase.go`)

**Types/Interfaces**:

- Use PascalCase for all type names
- Interface names should end with "er" (e.g., `ProductRepository`, `ProductUsecase`)
- Concrete implementations should match interface name without "er" prefix where possible (e.g., `productRepository`)

**Variables/Functions**:

- Use camelCase
- Avoid single-letter names except for short-lived loop variables
- Be descriptive: `productRepo` not `pr` or `repo`

**Constants**:

- Use PascalCase for exported constants
- Use camelCase for unexported constants

### Types

**Structs**:

```go
type Product struct {
    ID          uuid.UUID
    Name        string
    Description string
    CategoryID  uuid.UUID
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

**DTOs**:

- Request DTOs: `ProductReq`, `CreateProductRequest`
- Response DTOs: `ProductResp`, `ProductResponse`
- Use JSON tags for serialization

**Interfaces** (define in domain layer):

```go
type ProductRepository interface {
    Create(ctx context.Context, p ProductInput) error
    Fetch(ctx context.Context, limit, offset int) ([]Product, error)
    // ...
}
```

### Error Handling

- Return errors explicitly, don't use global error variables
- Wrap errors with context: `fmt.Errorf("failed to connect to database: %w", err)`
- In handlers, check errors and return appropriate HTTP status codes
- Use `errors.New()` for simple errors or custom error types for complex scenarios

```go
// Repository
if err != nil {
    return errors.New(err.Error())
}

// Handler
if err := h.usecase.CreateProduct(ctx, input); err != nil {
    c.JSON(http.StatusInternalServerError, dto.ErrorResp{
        Status:  http.StatusInternalServerError,
        Message: "failed to create product",
    })
    return
}
```

### Context Usage

- Always pass `context.Context` as first parameter to repository and usecase methods
- Extract context from request in handlers: `ctx := c.Request.Context()`
- Use context for cancellation and timeouts

### Dependency Injection

- Pass dependencies through constructor functions
- Use interfaces for dependencies (facilitates testing)

```go
func NewProductHandler(u usecase.ProductUsecase) *ProductHandler {
    return &ProductHandler{usecase: u}
}
```

### Database (sqlc)

- Write SQL queries in `sql/queries/*.sql` files
- Name queries with descriptive names: `GetProductByID`, `CreateProduct`, etc.
- Run `sqlc generate` after modifying SQL
- Generated code lives in `internal/db/`

### Logging

- Use the shared logger from `pkg/logger`
- Configure logger in main.go before use

```go
var log = logging.MustGetLogger("api")
log.Infof("Server starting on %s", serverAddr)
```

### Configuration

- Use `config/config.go` for environment-based config
- Use struct tags for env vars: `env:"PORT" env-default:"8080"`
- Use `cleanenv` for reading config

### HTTP Response Patterns

Use the standardized response DTOs from `internal/delivery/dto/response.go`:

```go
// Success responses
c.JSON(http.StatusOK, dto.Response{
    Status:  http.StatusOK,
    Message: "Success",
    Data:    result,
})

// Paginated responses
c.JSON(http.StatusOK, dto.PaginatedResponse{
    Status:     http.StatusOK,
    Message:    "Success",
    Data:       items,
    Total:      total,
    Page:       page,
    Limit:      limit,
    TotalPages: total / page,
})

// Error responses
c.JSON(http.StatusBadRequest, dto.ErrorResp{
    Status:  http.StatusBadRequest,
    Message: "invalid request body",
})
```

### Best Practices

1. **Keep layers separate**: Handler -> Usecase -> Repository
2. **Use interfaces**: Define repository interfaces in domain layer
3. **Validate early**: Validate input in handlers before passing to usecases
4. **Close resources**: Use defer for database connections, files, etc.
5. **Graceful shutdown**: Handle SIGINT/SIGTERM for clean server shutdown
6. **Use UUIDs**: Use `github.com/google/uuid` for IDs
7. **Write tests**: Create `*_test.go` files in the same package

### Common Patterns

**Handler flow**:

1. Extract context from request
2. Parse and validate input
3. Call usecase
4. Return appropriate response

**Usecase flow**:

1. Accept domain types as input
2. Call repository methods
3. Return domain types or errors

**Repository flow**:

1. Accept domain types as input
2. Convert to database params
3. Execute generated SQL
4. Convert results to domain types
5. Return domain types or errors

## Golden Rule

**Always prefix commands with `rtk`**. If RTK has a dedicated filter, it uses it. If not, it passes through unchanged. This means RTK is always safe to use.

**Important**: Even in command chains with `&&`, use `rtk`:

```bash
# ❌ Wrong
git add . && git commit -m "msg" && git push

# ✅ Correct
rtk git add . && rtk git commit -m "msg" && rtk git push
```

## RTK Commands by Workflow

### Build & Compile (80-90% savings)

```bash
rtk cargo build         # Cargo build output
rtk cargo check         # Cargo check output
rtk cargo clippy        # Clippy warnings grouped by file (80%)
rtk tsc                 # TypeScript errors grouped by file/code (83%)
rtk lint                # ESLint/Biome violations grouped (84%)
rtk prettier --check    # Files needing format only (70%)
rtk next build          # Next.js build with route metrics (87%)
```

### Test (90-99% savings)

```bash
rtk cargo test          # Cargo test failures only (90%)
rtk vitest run          # Vitest failures only (99.5%)
rtk playwright test     # Playwright failures only (94%)
rtk test <cmd>          # Generic test wrapper - failures only
```

### Git (59-80% savings)

```bash
rtk git status          # Compact status
rtk git log             # Compact log (works with all git flags)
rtk git diff            # Compact diff (80%)
rtk git show            # Compact show (80%)
rtk git add             # Ultra-compact confirmations (59%)
rtk git commit          # Ultra-compact confirmations (59%)
rtk git push            # Ultra-compact confirmations
rtk git pull            # Ultra-compact confirmations
rtk git branch          # Compact branch list
rtk git fetch           # Compact fetch
rtk git stash           # Compact stash
rtk git worktree        # Compact worktree
```

Note: Git passthrough works for ALL subcommands, even those not explicitly listed.

### GitHub (26-87% savings)

```bash
rtk gh pr view <num>    # Compact PR view (87%)
rtk gh pr checks        # Compact PR checks (79%)
rtk gh run list         # Compact workflow runs (82%)
rtk gh issue list       # Compact issue list (80%)
rtk gh api              # Compact API responses (26%)
```

### JavaScript/TypeScript Tooling (70-90% savings)

```bash
rtk pnpm list           # Compact dependency tree (70%)
rtk pnpm outdated       # Compact outdated packages (80%)
rtk pnpm install        # Compact install output (90%)
rtk npm run <script>    # Compact npm script output
rtk npx <cmd>           # Compact npx command output
rtk prisma              # Prisma without ASCII art (88%)
```

### Files & Search (60-75% savings)

```bash
rtk ls <path>           # Tree format, compact (65%)
rtk read <file>         # Code reading with filtering (60%)
rtk grep <pattern>      # Search grouped by file (75%)
rtk find <pattern>      # Find grouped by directory (70%)
```

### Analysis & Debug (70-90% savings)

```bash
rtk err <cmd>           # Filter errors only from any command
rtk log <file>          # Deduplicated logs with counts
rtk json <file>         # JSON structure without values
rtk deps                # Dependency overview
rtk env                 # Environment variables compact
rtk summary <cmd>       # Smart summary of command output
rtk diff                # Ultra-compact diffs
```

### Infrastructure (85% savings)

```bash
rtk docker ps           # Compact container list
rtk docker images       # Compact image list
rtk docker logs <c>     # Deduplicated logs
rtk kubectl get         # Compact resource list
rtk kubectl logs        # Deduplicated pod logs
```

### Network (65-70% savings)

```bash
rtk curl <url>          # Compact HTTP responses (70%)
rtk wget <url>          # Compact download output (65%)
```

### Meta Commands

```bash
rtk gain                # View token savings statistics
rtk gain --history      # View command history with savings
rtk discover            # Analyze Claude Code sessions for missed RTK usage
rtk proxy <cmd>         # Run command without filtering (for debugging)
rtk init                # Add RTK instructions to CLAUDE.md
rtk init --global       # Add RTK to ~/.claude/CLAUDE.md
```
