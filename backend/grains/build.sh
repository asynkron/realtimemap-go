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
# go install github.com/gogo/protobuf/protoc-gen-gogoslick@latest
# go install github.com/AsynkronIT/protoactor-go/protobuf/protoc-gen-gograinv2@latest
#
#
# (make sure go installation path is on your PATH)

protoc -I=. -I=$GOPATH/src --gogoslick_out=. protos.proto
protoc -I=. -I=$GOPATH/src --gograinv2_out=. protos.proto