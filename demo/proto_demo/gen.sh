go get -u google.golang.org/protobuf/cmd/protoc-gen-go@latest

go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

go mod tidy

protoc --go_out=../ --go-grpc_out=../ --gorm_out=../ ./*.proto