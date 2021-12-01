FROM golang:1.16 as builder

LABEL MAINTAINER="justxuewei"

WORKDIR /dubbogo

COPY ./dubbogo-client ./
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
