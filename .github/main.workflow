workflow "New workflow" {
  on = "push"
  resolves = ["gofmt"]
}

action "gofmt" {
  uses = "martinxxD/go-github-actions/go@master"
  secrets = ["GITHUB_TOKEN"]
}
