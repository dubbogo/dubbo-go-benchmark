FROM golang:1.16 as builder

WORKDIR /dubbogo

COPY ./ ./

RUN go mod init "github.com/dubbogo/dubbo-go-benchmark/apasvc-server" \
    && go mod tidy \
    # TODO(justxuewei): remove the below two lines if  \
    #   feat/adasvc branch was merged
    && go get dubbo.apache.org/dubbo-go/v3@feat-adasvc \
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
