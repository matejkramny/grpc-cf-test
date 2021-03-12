PROTO_INCLUDES=-I=./ \
		-I=/usr/local/include/google \
		-I=/usr/local/include \
		-I=$(GOPATH)/src

.PHONY: grpc server
all: grpc server

grpc:
	@protoc ${PROTO_INCLUDES} --go-grpc_out=./ --go_out=./ ./*.proto

server:
	@mkdir -p ./build
	go build -o ./build/server ./server