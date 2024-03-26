FROM golang:1.21.6

WORKDIR /go/src/github.com/gihyodocker/taskapp

# buildに必要なファイル群をコピー
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY go.mod .
COPY go.sum .
COPY Makefile .

# ビルドに関連する処理
RUN make mod
RUN make vendor
RUN make build-api

# APIサーバの実行ファイルをENTRYPOINTに設定
ENTRYPOINT ["./bin/api"]
