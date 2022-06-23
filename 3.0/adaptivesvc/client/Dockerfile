FROM golang:1.16 as builder

WORKDIR /dubbogo

COPY ./ ./

RUN go mod init "github.com/dubbogo/dubbo-go-benchmark/adasvc-client" \
    && go mod tidy \
    && go build -a -o dubbogo-client

FROM golang:1.16

WORKDIR /dubbogo

COPY --from=builder /dubbogo/dubbogo-client ./
COPY dubbogo.yml ./conf/dubbogo.yaml

ENV DUBBO_GO_CONFIG_PATH="./conf/dubbogo.yaml"
ENV TPS=1000
ENV PARALLEL=100
ENV DURATION="10m"
ENV FUNC_NAME="FIBONACCI"
ENV FIBONACCI_N=40
ENV FIBONACCI_WORKER_NUM=6

RUN chmod +x ./dubbogo-client

EXPOSE 20001

ENTRYPOINT ["./dubbogo-client"]
