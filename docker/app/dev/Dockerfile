# 軽量なalpineイメージを使う
FROM golang:1.14-alpine

# 必要な環境変数をダミー入力。これらの値はdocker-compose.ymlのenvironmentで上書きされる。
ENV POSTGRES_NAME dummydb
ENV POSTGRES_USER dummyuser
ENV POSTGRES_PASSWORD dummypassword
ENV POSTGRES_HOST dummyhost
ENV POSTGRES_PORT 1234
ENV JWT_SECRET_KEY dummyjwtkey

# Goビルドに必要なパッケージ群をインストールする
RUN apk add --no-cache git make musl-dev gcc

# タイムゾーンをJSTに設定する
RUN apk --no-cache add tzdata && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

# imports管理パッケージgoimportsをインストール。hack/lint.shで使用している。
RUN go get -u golang.org/x/tools/cmd/goimports

# linterパッケージgolintをインストール。
RUN go get -u golang.org/x/lint/golint

# 自動リビルドパッケージgodoをインストール。
RUN go get -u gopkg.in/godo.v2/cmd/godo

# 依存パッケージのインストール
WORKDIR /go/src/github.com/gold-kou/go-housework
COPY go.mod go.sum ./
RUN go mod download

# Goアプリケーションのビルド
ADD . /go/src/github.com/gold-kou/go-housework
RUN make install