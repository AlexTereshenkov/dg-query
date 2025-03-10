build-go:
	gofmt -w .
	goimports -w .
	staticcheck ./...
	go build
	go install
	go test ./...
	dg-query

run-go:
	go build
	go run main.go dependencies --dg="examples/dg-real.json" foo.py spam.py

build-bazel:
	.buildkite/build-bazel.sh
