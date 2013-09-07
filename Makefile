test:
	@go get launchpad.net/gocheck
	@go test -i ./dsl
	go test ./dsl
	@go run gspec.go specs
