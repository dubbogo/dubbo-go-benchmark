FROM golang:1.16 as builder

WORKDIR /dubbogo

COPY ./ ./

RUN go mod init "github.com/dubbogo/dubbo-go-benchmark/apasvc-client" \
    && go mod tidy \
    # TODO(justxuewei): remove the below two lines if  \
    #   feat/adasvc branch was merged
    && go get dubbo.apache.org/dubbo-go/v3@feat-adasvc \
    && go mod tidy \
    && go build -a -o dubbogo-client

FROM golang:1.16

LABEL MAINTAINER="justxuewei"

WORKDIR /dubbogo

COPY --from=builder /dubbogo/dubbogo-client ./
COPY dubbogo.yml ./

ENV TPS=1000
ENV PARALLEL=100
ENV DURATION="10m"
ENV FUNC_NAME="FIBONACCI"
ENV FIBONACCI_N=40
ENV FIBONACCI_WORKER_NUM=6

RUN chmod +x ./dubbogo-client

EXPOSE 20001

ENTRYPOINT ["./dubbogo-client"]
