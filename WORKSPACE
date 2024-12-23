load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

# to support spawning actions in a Docker container
register_execution_platforms(
    "//platforms:docker",
)

http_archive(
    name = "com_github_sluongng_nogo_analyzer",
    sha256 = "0dc6b5e86094d081e05bcd0c3e41fc275a2398c64e545376166139412181f150",
    strip_prefix = "nogo-analyzer-0.0.3",
    urls = [
        "https://github.com/sluongng/nogo-analyzer/archive/refs/tags/v0.0.3.tar.gz",
    ],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "d93ef02f1e72c82d8bb3d5169519b36167b33cf68c252525e3b9d3d5dd143de7",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.49.0/rules_go-v0.49.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.49.0/rules_go-v0.49.0.zip",
    ],
)

# gazelle setup
http_archive(
    name = "bazel_gazelle",
    sha256 = "d76bf7a60fd8b050444090dfa2837a4eaf9829e1165618ee35dceca5cbdf58d5",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.37.0/bazel-gazelle-v0.37.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(
    nogo = "@//:nogo", 
    version = "1.22.5",
)

# This should be called before the `gazelle_dependencies()` as it seems that
# some of the dependencies (in the shared Bazel dependences space) are overriden
load("@com_github_sluongng_nogo_analyzer//staticcheck:deps.bzl", "staticcheck")
staticcheck()

gazelle_dependencies()

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

http_archive(
    name = "lcov",
    urls = ["https://github.com/linux-test-project/lcov/releases/download/v2.1/lcov-2.1.tar.gz"],
    strip_prefix = "lcov-2.1",
    sha256 = "4d01d9f551a3f0e868ce84742fb60aac4407e3fc1622635a07e29d70e38f1faf",
    build_file_content = """
package(default_visibility = ["//visibility:public"])

filegroup(
    name = "genhtml",
    srcs = glob([
        "bin/genhtml",
    ]),
)

filegroup(
    name = "bin",
    srcs = glob([
        "bin/*",
    ]),
)
    """
)

# Buildifier setup
http_file(
    name = "buildifier",
    urls = ["https://github.com/bazelbuild/buildtools/releases/download/v7.3.1/buildifier-linux-amd64"],
    sha256 = "5474cc5128a74e806783d54081f581662c4be8ae65022f557e9281ed5dc88009",
    downloaded_file_path = "buildifier",
    executable = True,
)
