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
	bazel run //:gazelle
	bazel run //:buildifier
	bazel run @rules_go//go -- fmt ./...
	bazel run @rules_go//go -- mod tidy
	bazel build //...
	bazel test //...
