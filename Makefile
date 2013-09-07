test:
	@go test -i ./dsl
	go test ./dsl
	@go run gspec.go specs
