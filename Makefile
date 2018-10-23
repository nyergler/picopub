setup:
	go install golang.org/x/tools/cmd/stringer
	go install github.com/campoy/jsonenums

deps:
	dep ensure

generate:
	go generate ./...

test: generate
	go test -v ./...