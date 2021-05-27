.PHONY: stub proto build

stub:
	@echo "It's a stub"

proto: 
	export PATH=$$PATH:$$GOPATH/bin;\
	protoc -I $$GOPATH/src/NewPhoto  --go_out=. --go-grpc_out=. logic/proto/api.proto

build: proto
	go build main.go