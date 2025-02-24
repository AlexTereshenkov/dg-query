bazel run //:gazelle
bazel run //:buildifier
bazel run //:buildifier-check

bazel mod tidy

# this would cause nogo analyzer to fail when running `bazel build //...`
# obj := false
# res := fmt.Sprintf("%d", obj)
# fmt.Println(res)
# go run main.go subgraph --dg="examples/dg-linked-list.json" --root="7"

bazel test //...
