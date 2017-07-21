load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_prefix", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["nacl.go"],
    visibility = ["//visibility:public"],
    deps = ["//randombytes:go_default_library"],
)

go_prefix("github.com/kevinburke/nacl")

go_test(
    name = "go_default_test",
    srcs = ["nacl_test.go"],
    timeout = "short",
    library = ":go_default_library",
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    timeout = "short",
    deps = ["//:go_default_library"],
)
