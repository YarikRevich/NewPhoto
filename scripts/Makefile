.PHONY: stub proto build

stub:
	@echo "It's a stub"

proto: 
	export PATH=$$PATH:$$GOPATH/bin;\
	protoc -I $$GOPATH/src/NewPhoto  --go_out=$$GOPATH/src/NewPhoto --go-grpc_out=$$GOPATH/src/NewPhoto $$GOPATH/src/NewPhoto/logic/proto/api.proto

build: proto
	go build main.go
