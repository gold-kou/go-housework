# ビルド用
## 軽量なalpineイメージを使う
FROM golang:1.14-alpine AS golang-build

## Goビルドに必要なアプリケーションをインストールする
RUN apk add --no-cache git make musl-dev gcc

## 依存パッケージのインストール
WORKDIR /go/src/github.com/gold-kou/go-housework
COPY go.mod go.sum ./
RUN go mod download

## GitHubリポジトリのファイルをコピー
ADD . /go/src/github.com/gold-kou/go-housework

## Goアプリケーションのビルド
RUN make cross-install

## 環境変数を読み込んでビルド
ENV GOPATH /go
ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0
WORKDIR /tmp
RUN git clone https://github.com/ProjectThor/get-ssm-params \
  && cd get-ssm-params \
  && go get -d
RUN  cd get-ssm-params ; go build

# 実行用(マルチステージビルド)
## ベースイメージの再指定
FROM alpine

## タイムゾーンをJSTに設定する
RUN apk --no-cache add tzdata && cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

## 必要な環境変数をダミー入力
ENV POSTGRES_NAME dummydb
ENV POSTGRES_USER dummyuser
ENV POSTGRES_PASSWORD dummypassword
ENV POSTGRES_HOST dummyhost
ENV POSTGRES_PORT 1234
ENV JWT_SECRET_KEY dummyjwtkey

## Goアプリケーションをコピー
COPY --from=golang-build /go/bin/zozoerp-* /

## SQLをコピー
COPY --from=golang-build /go/src/github.com/gold-kou/go-housework/app/sql /app/sql

## docker-entrypoint.shをコピー
COPY --from=golang-build /go/src/github.com/gold-kou/go-housework/docker-entrypoint.sh /

## SSMライブラリをコピー
COPY --from=golang-build /tmp/get-ssm-params/get-ssm-params /

## 公開ポート設定
EXPOSE 8080

## このイメージが起動したらdocker-entrypoint.shを実行
ENTRYPOINT ["./docker-entrypoint.sh"]

## このイメージが起動したらアプリケーション起動
CMD ["./run-server"]
