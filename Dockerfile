FROM ubuntu:18.04
LABEL maintainer="saquibulhassan6@gmail.com"

RUN apt-get update -y
RUN apt-get install -y squid
RUN apt-get install -y apache2-utils wget
RUN apt-get install -y sqlite3 build-essential
ENV SQUID_VERSION=3.5.27 \
    SQUID_CACHE_DIR=/var/spool/squid \
    SQUID_LOG_DIR=/var/log/squid \
    SQUID_USER=proxy
COPY squid.main.conf /etc/squid/squid.conf
COPY entrypoint.sh /sbin/entrypoint.sh
RUN chmod 755 /sbin/entrypoint.sh
RUN service squid start
EXPOSE 3128/tcp
ENTRYPOINT ["/sbin/entrypoint.sh"]

RUN wget https://dl.google.com/go/go1.13.7.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.13.7.linux-amd64.tar.gz
RUN rm go1.13.7.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"


RUN mkdir /app

ADD . /app

WORKDIR /app
RUN go get -u github.com/gorilla/mux
RUN go get -u github.com/mattn/go-sqlite3

RUN go build -o main .
CMD ["/app/main"]