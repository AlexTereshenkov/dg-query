bazel run //:gazelle
bazel run //:buildifier
bazel run //:buildifier.check

bazel mod tidy