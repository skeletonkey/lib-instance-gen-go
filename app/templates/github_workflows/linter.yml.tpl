name: golangci-lint

on:
  pull_request:
    branches:
      - "*"

jobs:
  golangci:
    name: lint
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: {{ .GoVersion }}
          cache: false

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
        with:
          version: v1.54
