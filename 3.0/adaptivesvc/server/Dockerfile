FROM golang:1.16 as builder

LABEL MAINTAINER="justxuewei"

WORKDIR /dubbogo

COPY ./dubbogo-server ./
COPY dubbogo.yml ./

RUN chmod +x ./dubbogo-server

EXPOSE 20001

ENTRYPOINT ["./dubbogo-server"]
