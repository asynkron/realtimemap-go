# To run these commands 
# 1. Install protoc compiler
#
# (on Windows)
# choco install protoc
#
# (on Ubuntu)
# apt install protobuf-compiler
#
# 2. Install plugins
# 
# go install github.com/golang/protobuf/protoc-gen-go@latest
# go install github.com/asynkron/protoactor-go/protobuf/protoc-gen-gograinv2@latest
#
#
# (make sure go installation path is on your PATH)

protoc -I=. -I=$GOPATH/src --go_opt=paths=source_relative --go_out=. protos.proto
protoc -I=. -I=$GOPATH/src --gograinv2_out=. protos.proto