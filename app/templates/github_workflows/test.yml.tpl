name: tests

on:
  pull_request:
    branches:
      - "*"

jobs:
  build:
    name: tests
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: "{{ .GoVersion }}"

      - name: Prepare dependencies
        run: |
          go mod download

      - name: Test code
        run: |
          go test -race ./...
