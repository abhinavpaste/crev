# crev

A lightweight code review tool built in Go. Share code snippets and leave inline line-level comments — no bloat, no framework, just Go's standard library and PostgreSQL.

## Tech Stack

- **Language:** Go (net/http — no frameworks)
- **Database:** PostgreSQL
- **Auth:** JWT (HS256)
- **Password Hashing:** bcrypt

## Project Structure

```
crev/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── auth/
│   │   └── auth.go          # JWT generation/validation, bcrypt hashing
│   ├── db/
│   │   ├── db.go            # PostgreSQL connection pool
│   │   └── schema.sql       # Table definitions
│   ├── handler/
│   │   ├── auth.go          # Register, Login handlers
│   │   ├── snippet.go       # Create, Get snippet handlers
│   │   └── comment.go       # Create, Get comment handlers
│   ├── middleware/
│   │   └── auth.go          # JWT auth middleware
│   └── models/
│       └── models.go        # User, Snippet, Comment structs
├── .env
├── go.mod
└── go.sum
```

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL

### 1. Clone the repo

```bash
git clone https://github.com/abhinavpaste/crev.git
cd crev
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Create the database

```bash
psql -U postgres -c "CREATE DATABASE crev;"
psql -U postgres -d crev -f internal/db/schema.sql
```

### 4. Configure environment

Create a `.env` file in the project root:

```
DATABASE_URL=postgres://postgres:yourpassword@localhost:5432/crev?sslmode=disable
JWT_SECRET=your_secret_key
```

### 5. Run

```bash
go run ./cmd/server/main.go
```

Server starts on `http://localhost:8080`.

## API Reference

### Auth

#### Register
```
POST /register
Content-Type: application/json

{
  "username": "abhinav",
  "email": "abhinav@example.com",
  "password": "secret123"
}
```

**Response:** `201 Created` — user object (password hash excluded)

---

#### Login
```
POST /login
Content-Type: application/json

{
  "email": "abhinav@example.com",
  "password": "secret123"
}
```

**Response:** `200 OK`
```json
{
  "token": "<jwt>"
}
```

---

### Snippets

#### Create Snippet
```
POST /snippets
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Hello World",
  "content": "fmt.Println(\"Hello, World!\")",
  "language": "go"
}
```

**Response:** `201 Created` — snippet object with ID

---

#### Get Snippet
```
GET /snippets/:id
```

**Response:** `200 OK` — snippet object (public, no auth required)

---

### Comments

#### Add Comment
```
POST /snippets/:id/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "line_number": 3,
  "content": "this could use a better variable name"
}
```

**Response:** `201 Created` — comment object

---

#### Get Comments
```
GET /snippets/:id/comments
```

**Response:** `200 OK` — array of comments ordered by line number (public, no auth required)

---

## Schema

```sql
users    → id, username, email, password_hash, created_at
snippets → id, user_id, title, content, language, created_at
comments → id, snippet_id, user_id, line_number, content, created_at
```

## Design Decisions

- **No framework** — built entirely on `net/http` to keep the binary lean and dependencies minimal
- **JWT stateless auth** — tokens expire in 72 hours, no session store needed
- **bcrypt cost 14** — deliberately slow to mitigate brute force on password hashes
- **UUID primary keys** — avoids sequential ID enumeration attacks
- **Public snippet reads** — snippets and comments are readable without auth, write operations require a valid JWT
