workflow "New workflow" {
  on = "push"
  resolves = ["gofmt"]
}

action "gofmt" {
  uses = "martinxxD/go-github-actions/fmt@master"
  secrets = ["GITHUB_TOKEN"]
}
