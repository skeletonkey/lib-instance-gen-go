name: Go

on:
  push:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v6

    - name: Set up Go
      uses: actions/setup-go@v6.1

    - name: Download dependencies
      run: go mod tidy

    - name: Test
      run: go test -v ./...