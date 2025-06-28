.PHONY: proto run-server run-client test clean tidy

proto:
	@echo "Generating gRPC files..."
		protoc \
			--go_out=gen/chat --go_opt=paths=source_relative \
			--go-grpc_out=gen/chat --go-grpc_opt=paths=source_relative \
			proto/chat.proto

run-server:
	go run server/main.go

run-client:
	go run client/cmd/main.go

test:
	go test ./...

clean:
	rm -rf gen/chat/proto

tidy:
	cd proto && go mod tidy
	cd server && go mod tidy
	cd client && go mod tidy