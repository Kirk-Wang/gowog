#!/bin/bash
echo "generating proto"

protoc -I . -I $GOPATH/src --go_out=Message_proto/. message.proto
protoc -I . --js_out=import_style=commonjs,binary:$GOPATH/src/gowog-cloud/client/src/states message.proto
protoc -I . --python_out=../ai/ message.proto
