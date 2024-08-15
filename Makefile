build-go:
	go build
	go install
	dg-query

run-go:
	go build
	go run main.go dependencies --dg="examples/dg-real.json" foo.py spam.py

build-bazel:
	bazel run //:dg-query
