# graphql-gen-go

Code generation for GraphQL clients in Golang.

Features:
- Generate JSON-able structs for all types in provided schema.
- Generate interfaces representing views provided by queries.
- TODO: generate full client
- TODO: generate helpers for httptest

Example:
See `internal/test/swapi`.
- Change/create a query, regenerate with `go generate ./...`
- Check out `api_test.go` for sample usage.
