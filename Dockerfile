# 1段目: Go のプログラムを組み立てる(ビルド)
FROM golang:1.27 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o server .

# 2段目: できあがったプログラムだけを、小さな箱に入れる
FROM gcr.io/distroless/static-debian12
COPY --from=build /app/server /server
CMD ["/server"]