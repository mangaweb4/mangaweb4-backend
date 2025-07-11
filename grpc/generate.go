package grpc

//go:generate protoc --go_out=. --proto_path=./schema --go_opt=paths=source_relative manga.proto
//go:generate protoc --go_out=. --proto_path=./schema --go_opt=paths=source_relative history.proto
//go:generate protoc --go_out=. --proto_path=./schema --go_opt=paths=source_relative tag.proto
//go:generate protoc --go_out=. --proto_path=./schema --go_opt=paths=source_relative types.proto
//go:generate protoc --go_out=. --proto_path=./schema --go_opt=paths=source_relative maintenance.proto
