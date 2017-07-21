git_repository(
    name = "io_bazel_rules_go",
    remote = "https://github.com/bazelbuild/rules_go.git",
    tag = "0.5.2",
)
load("@io_bazel_rules_go//go:def.bzl", "go_prefix", "go_repositories", "go_repository")

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    commit = "6914964337150723782436d56b3f21610a74ce7b",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    commit = "18107e6c56edb2d51f965f7d68e59404f0daee54",
)

go_repositories()
