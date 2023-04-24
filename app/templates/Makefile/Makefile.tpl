build:
	go fmt ./...
	go vet ./...
	go build -o bin/{{ .BinaryName }} app/*.go

install:
	cp bin/{{ .BinaryName }} /usr/local/sbin/{{ .BinaryName }}

golib-latest:
	go get -u github.com/skeletonkey/lib-instance-gen-go@latest

	go mod tidy

app-init:
	go generate
