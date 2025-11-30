name: golangci-lint
on:
  pull_request:

permissions:
  contents: read
  # Optional: allow read access to pull request. Use with `only-new-issues` option.
  # pull-requests: read

jobs:
  golangci:
    name: lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: actions/setup-go@v6.1
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest