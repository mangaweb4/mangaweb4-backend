package grpc

//go:generate protoc -I$PROTOC_INCLUDE --go_out=. --go-grpc_out=. --proto_path=./schema --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative manga.proto history.proto tag.proto types.proto maintenance.proto
