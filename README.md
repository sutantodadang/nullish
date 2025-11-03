# Nullish

[![Go Reference](https://pkg.go.dev/badge/github.com/sutantodadang/nullish.svg)](https://pkg.go.dev/github.com/sutantodadang/nullish)
[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](https://choosealicense.com/licenses/mit/)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

A high-performance, minimal-dependency Go package for handling nullable types with PostgreSQL and JSON serialization support.

## Features

- High Performance - Sub-nanosecond operations for primitive types, zero allocations
- Minimal Dependencies - Only 3 required packages (go-json, uuid, ulid)
- Fully Tested - 77.5% code coverage, 87 test cases, 34 benchmarks
- SQL Driver Compatible - Implements database/sql/driver.Valuer and sql.Scanner
- JSON Support - Custom JSON marshaling/unmarshaling for all types
- Production Ready - Battle-tested in production environments

## Supported Types

| Type       | Description        | Use Case                    |
| ---------- | ------------------ | --------------------------- |
| NullString | Nullable string    | VARCHAR, TEXT columns       |
| NullInt    | Nullable integer   | INT, BIGINT columns         |
| NullFloat  | Nullable float64   | FLOAT, DOUBLE columns       |
| NullBool   | Nullable boolean   | BOOLEAN columns             |
| NullTime   | Nullable time.Time | TIMESTAMP columns           |
| NullUUID   | Nullable UUID      | UUID columns                |
| NullULID   | Nullable ULID      | Sortable unique identifiers |
| NullJSON   | Raw JSON           | JSONB columns               |
| NullObj    | JSON object        | map[string]interface{}      |
| NullArr    | JSON array         | []interface{}               |
| NullArrObj | Array of objects   | []map[string]interface{}    |

## Installation

```bash
go get -u github.com/sutantodadang/nullish
```

## Quick Start

### Basic Usage

```go
import "github.com/sutantodadang/nullish"

type User struct {
    ID        int                    `json:"id"`
    Name      nullish.NullString     `json:"name"`
    Age       nullish.NullInt        `json:"age"`
    Email     nullish.NullString     `json:"email"`
    Active    nullish.NullBool       `json:"active"`
    CreatedAt nullish.NullTime       `json:"created_at"`
    Metadata  nullish.NullObj        `json:"metadata"`
}

// Create with constructor
user := User{
    Name:  nullish.NewNullString("John Doe", true),
    Age:   nullish.NewNullInt(30, true),
    Email: nullish.NewNullString("", false), // null email
}
```

### Database Operations

```go
import (
    "database/sql"
    "github.com/sutantodadang/nullish"
)

// Scanning from database
var user User
err := db.QueryRow("SELECT name, age, email FROM users WHERE id = $1", 1).
    Scan(&user.Name, &user.Age, &user.Email)

// Inserting to database
_, err = db.Exec(
    "INSERT INTO users (name, age, email) VALUES ($1, $2, $3)",
    user.Name, user.Age, user.Email,
)
```

### JSON Serialization

```go
import "encoding/json"

// Marshal to JSON
user := User{
    Name:  nullish.NewNullString("Alice", true),
    Age:   nullish.NewNullInt(25, true),
    Email: nullish.NewNullString("", false),
}

jsonData, _ := json.Marshal(user)
// Output: {"id":0,"name":"Alice","age":25,"email":null,...}

// Unmarshal from JSON
var decoded User
json.Unmarshal(jsonData, &decoded)
```

### Complex Types

```go
// NullObj - for JSON objects
metadata := nullish.NewNullObj(map[string]interface{}{
    "role": "admin",
    "permissions": []string{"read", "write"},
}, true)

// NullArr - for arrays
tags := nullish.NewNullArr([]interface{}{"golang", "backend", "api"}, true)

// NullArrObj - for array of objects
addresses := nullish.NewNullArrObj([]map[string]interface{}{
    {"street": "123 Main St", "city": "NYC"},
    {"street": "456 Oak Ave", "city": "LA"},
}, true)
```

### UUID & ULID

```go
import (
    "github.com/google/uuid"
    "github.com/oklog/ulid/v2"
)

// UUID
id := uuid.New()
userID := nullish.NewNullUUID(id, true)

// ULID (sortable, timestamp-based)
entropy := ulid.DefaultEntropy()
ulidValue := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
trackingID := nullish.NewNullULID(ulidValue, true)
```

## Performance

Benchmark results on AMD Ryzen 5 7500F:

| Operation             | Time/op | Allocations |
| --------------------- | ------- | ----------- |
| NullString Value      | 0.21 ns | 0           |
| NullInt Scan          | 1.87 ns | 0           |
| NullBool Value        | 0.21 ns | 0           |
| NullFloat Scan        | 1.75 ns | 0           |
| NullTime MarshalJSON  | 124 ns  | 112 B       |
| NullUUID Scan (bytes) | 19 ns   | 0           |
| NullULID Value        | 6.8 ns  | 0           |
| NullObj UnmarshalJSON | 203 ns  | 208 B       |
| NullArr MarshalJSON   | 75 ns   | 32 B        |
| NullArrObj Scan       | 0.42 ns | 0           |

## Architecture

All types implement:

- database/sql/driver.Valuer - for database writes
- database/sql.Scanner - for database reads
- json.Marshaler - for JSON encoding
- json.Unmarshaler - for JSON decoding

Each type has two fields:

- The actual value (e.g., String, Int, Time)
- Valid bool - indicates if the value is non-null

## Testing

Run all tests:

```bash
go test -v
```

Run with coverage:

```bash
go test -cover
```

Run benchmarks:

```bash
go test -bench . -benchmem
```

## API Reference

### Constructors

All types provide constructor functions:

```go
NewNullString(str string, valid bool) NullString
NewNullInt(integer int, valid bool) NullInt
NewNullFloat(float float64, valid bool) NullFloat
NewNullBool(boolean bool, valid bool) NullBool
NewNullTime(time time.Time, valid bool) NullTime
NewNullUUID(uuid uuid.UUID, valid bool) NullUUID
NewNullULID(ulid ulid.ULID, valid bool) NullULID
NewNullJSON(json json.RawMessage, valid bool) NullJSON
NewNullObj(object map[string]interface{}, valid bool) NullObj
NewNullArr(array []interface{}, valid bool) NullArr
NewNullArrObj(arrayObject []map[string]interface{}, valid bool) NullArrObj
```

### Methods

All types implement:

- Value() (driver.Value, error) - Convert to database value
- Scan(value interface{}) error - Read from database
- MarshalJSON() ([]byte, error) - Convert to JSON
- UnmarshalJSON(data []byte) error - Parse from JSON

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

- [@sutantodadang](https://www.github.com/sutantodadang)

## Acknowledgments

- Uses [goccy/go-json](https://github.com/goccy/go-json) for high-performance JSON operations
- Uses [google/uuid](https://github.com/google/uuid) for UUID support
- Uses [oklog/ulid](https://github.com/oklog/ulid) for ULID support
