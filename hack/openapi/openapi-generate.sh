#!/usr/bin/env bash

GENERATE_PATH='openapi/openapi-generated'

# code generate
docker run -v ${PWD}:/local openapitools/openapi-generator-cli:v4.2.3 generate \
  -i /local/openapi/openapi.yaml \
  -g go-server \
  -p withXml=true \
  --output /local/${GENERATE_PATH}

# generateされたrouters.goをいったんgeneratedファイルとする。
# このgeneratedファイルをベースに手動で認証系の処理を追記し、app/server/routers.goとする。（自動生成の上書きを防ぐため）
mv ${GENERATE_PATH}/go/routers.go ${GENERATE_PATH}/routers.go.generated

# 自動生成されたスキーマモデルをschemamodelディレクトリに移動
mv ${GENERATE_PATH}/go/model*.go app/model/schemamodel

# schemamodelディレクトリ配下ファイルのpackage名をschemamodelに変更
sed -i '' -e 's/package openapi/package schemamodel/g' app/model/schemamodel/*.go
