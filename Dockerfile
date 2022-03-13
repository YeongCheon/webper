FROM ubuntu:20.04
LABEL maintainer="kyc1682@gmail.com"

RUN apt update
RUN apt install -y ca-certificates
RUN update-ca-certificates

COPY ./.bin ./.bin
COPY webper .

ENTRYPOINT ["./webper"]
