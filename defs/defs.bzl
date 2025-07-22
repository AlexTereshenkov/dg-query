"""Macros and shared definition."""

load("@rules_go//go:def.bzl", "go_test")
load("//:defs/constants.bzl", "EXAMPLES_DG_JSON", "TRANSITIVE_REDUCTION_DG_JSON")

def _custom_go_test_impl(name, visibility, srcs, data, deps, tags):
    go_test(
        name = name,
        data = (data or []) + [EXAMPLES_DG_JSON, TRANSITIVE_REDUCTION_DG_JSON],
        deps = (deps or []) + ["//cmd", "@com_github_stretchr_testify//assert", "@com_github_spf13_cast//:cast"],
        srcs = srcs,
        # running `bazel test --config=windows //tests:all` on Windows would skip these tests
        tags = tags,
        target_compatible_with = ["@platforms//os:linux"],
        visibility = visibility,
    )

custom_go_test = macro(
    attrs = {
        "data": attr.label_list(configurable = False),
        "deps": attr.label_list(configurable = False),
        "srcs": attr.label_list(configurable = False, mandatory = True),
        "tags": attr.string_list(configurable = False, mandatory = True),
    },
    implementation = _custom_go_test_impl,
)
