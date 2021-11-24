FROM golang:1.16 as builder

WORKDIR /dubbogo

COPY ./ ./

RUN go mod init "github.com/dubbogo/dubbo-go-benchmark" \
    && go mod tidy \
    && go build -a -o dubbogo-server

FROM golang:1.16

LABEL MAINTAINER="justxuewei"

WORKDIR /dubbogo

COPY --from=builder /dubbogo/dubbogo-server ./
COPY dubbogo.yml ./

RUN chmod +x ./dubbogo-server

EXPOSE 20001

ENTRYPOINT ["./dubbogo-server"]
