build-go:
	go build
	go install
	dg-query

build-bazel:
	bazel run //:dg-query