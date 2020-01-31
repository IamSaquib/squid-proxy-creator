FROM ubuntu:18.04
LABEL maintainer="Saquib Ul Hassan <saquibulhassan6@gmail.com>"

RUN apt-get update -y
RUN apt-get install -y squid
RUN apt-get install -y apache2-utils wget

RUN wget https://dl.google.com/go/go1.13.7.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.13.7.linux-amd64.tar.gz
RUN ls
RUN rm go1.13.7.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
RUN go version
# COPY squid.conf /etc/squid/squid.conf

# EXPOSE 3128/tcp

# RUN mkdir /app

# ADD main.go /app

# WORKDIR /app

# RUN go build -o main .

# CMD ["/app/main"]