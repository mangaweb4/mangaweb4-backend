package grpc

//go:generate protoc --go_out=. --proto_path=./schema --go_opt=paths=source_relative browse.proto
//go:generate protoc --go_out=. --proto_path=./schema --go_opt=paths=source_relative view.proto
