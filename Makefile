test:
	go get launchpad.net/gocheck
	go test -i .
	go test .
	go run gspec/gspec.go foo
	go run gspec/gspec.go specs
