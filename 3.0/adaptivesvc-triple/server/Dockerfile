FROM golang:1.16 as builder

WORKDIR /dubbogo

COPY ./ ./

RUN  go mod init "github.com/dubbogo/dubbo-go-benchmark/adasvc-server" \
    && go mod tidy \
    && go build -a -o dubbogo-server

FROM golang:1.16

WORKDIR /dubbogo

COPY --from=builder /dubbogo/dubbogo-server ./
COPY dubbogo.yml ./conf/dubbogo.yaml

ENV DUBBO_GO_CONFIG_PATH="./conf/dubbogo.yaml"

RUN chmod +x ./dubbogo-server

EXPOSE 20001

ENTRYPOINT ["./dubbogo-server"]
