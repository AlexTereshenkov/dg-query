#! /bin/bash -xue
echo --- 'Prepare environment'
sudo apt-get install wget jq -y
wget https://dl.google.com/go/go1.22.7.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.7.linux-amd64.tar.gz

echo --- 'Check Go installation'
export GOROOT="/usr/local/go"
# setting GOBIN is required for `go install` output location
export GOBIN="/usr/local/go/bin/"
export PATH="$GOROOT/bin:$PATH"
go version
go env

echo --- 'Static checks'
go install golang.org/x/tools/cmd/goimports@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
ls /usr/local/go/bin

if [ -n "$(gofmt -l -s .)" ]; then
    echo "Code needs to be formatted"
    exit 1
fi

if [ -n "$(goimports -l .)" ]; then
    echo "Imports need to be formatted"
    exit 1
fi

staticcheck ./...

echo --- 'Build project'
go build
mkdir dist
go test -test.v -coverprofile dist/coverage.out  ./...
go tool cover -html dist/coverage.out -o dist/coverage.html

echo --- 'Sanity check'
go run main.go deps --dg='examples/dg-real.json' foo.py spam.py
go run main.go rdeps --dg='examples/dg-real.json' eggs.py
go run main.go metrics --dg="examples/dg-real-transitive.json" --metric=deps-transitive | jq
