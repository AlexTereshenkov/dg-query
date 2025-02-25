"""Macros and shared definition."""

load("@rules_go//go:def.bzl", "go_test")

def custom_go_test(
        name,
        srcs,
        tag,
        data = ["examples/dg.json"],
        deps = ["//cmd", "@com_github_stretchr_testify//assert", "@com_github_spf13_cast//:cast"]):
    go_test(
        name = name,
        srcs = srcs,
        data = data,
        deps = deps,
        tags = [tag],
        target_compatible_with = ["@platforms//os:linux"],
    )
