http_archive(
    name = "io_bazel_rules_go",
    sha256 = "91fca9cf860a1476abdc185a5f675b641b60d3acf0596679a27b580af60bf19c",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.7.0/rules_go-0.7.0.tar.gz",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()

go_register_toolchains()

go_repository(
    name = "com_github_google_go_cmp",
    commit = "98232909528519e571b2e69fbe546b6ef35f5780",
    importpath = "github.com/google/go-cmp",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "bd6f299fb381e4c3393d1c4b1f0b94f5e77650c8",
    importpath = "golang.org/x/crypto",
)
