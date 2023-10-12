.DEFAULT_GOAL=build

build:
	go fmt ./...
	go vet ./...
	{{.BuildEnvArgs}}go build -o bin/{{ .BinaryName }} app/*.go

install:
	cp bin/{{ .BinaryName }} /usr/local/sbin/{{ .BinaryName }}

golib-latest:
	{{ .Dependencies }}go get -u github.com/skeletonkey/lib-core-go@latest
	go get -u github.com/skeletonkey/lib-instance-gen-go@latest

	go mod tidy

app-init:
	go generate
