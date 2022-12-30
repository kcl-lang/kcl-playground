FROM golang:1.17 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o kcl-playground


FROM kusionstack/kusion:v0.7.2

WORKDIR /app

COPY --from=builder /app/kcl-playground .

EXPOSE 2022


CMD ["./kcl-playground"]