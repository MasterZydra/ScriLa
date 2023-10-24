# Development

**Compile and run**  
`go run ./...`

**Build executable**  
`go build -o . ./...`

**Run all tests**  
`go test -v ./...`

**See test coverage**  
`go test -coverprofile=coverage.out ./...`  
`go tool cover -html=coverage.out`
