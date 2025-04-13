"""Macros and shared definition."""

load("@rules_go//go:def.bzl", "go_test")
load("//:defs/constants.bzl", "EXAMPLES_DG_JSON", "TRANSITIVE_REDUCTION_DG_JSON")

def custom_go_test(
        name,
        srcs,
        tag,
        data = [EXAMPLES_DG_JSON, TRANSITIVE_REDUCTION_DG_JSON],
        deps = ["//cmd", "@com_github_stretchr_testify//assert", "@com_github_spf13_cast//:cast"]):
    go_test(
        name = name,
        srcs = srcs,
        data = data,
        deps = deps,
        tags = [tag],
        # running `bazel test --config=windows //tests:all` on Windows would skip these tests
        target_compatible_with = ["@platforms//os:linux"],
    )
