"""Macros and shared definition."""

load("@io_bazel_rules_go//go:def.bzl", "go_test")

def custom_go_test(name, srcs, tag, data = ["examples/dg.json"], deps = ["//cmd", "@com_github_stretchr_testify//assert"]):
    go_test(
        name = name,
        srcs = srcs,
        data = data,
        deps = deps,
        tags = [tag],
        target_compatible_with = ["@platforms//os:linux"],
    )
