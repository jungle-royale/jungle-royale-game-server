#!/bin/bash

# 각 폴더 protoc 컴파일 스크립트

# 스크립트가 위치한 디렉토리 기준으로 작업
BASE_DIR=$(dirname "$0")

# server 폴더에서 Go 코드 생성
echo "Generating Go code in server folder..."
cd "$BASE_DIR" || exit
# pwd
protoc --go_out=. message.proto
if [ $? -ne 0 ]; then
    echo "Failed to generate Go code."
    exit 1
fi

# phaser-client 폴더에서 JavaScript 코드 생성
echo "Generating JavaScript code in phaser-client folder..."
cd "$BASE_DIR/phaser-client" || exit
pwd
protoc --js_out=import_style=commonjs,binary:. message.proto
if [ $? -ne 0 ]; then
    echo "Failed to generate JavaScript code."
    exit 1
fi

echo "Protobuf code generation completed successfully!"