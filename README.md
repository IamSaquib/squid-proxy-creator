# Squid Proxy Creator

A Squid Proxy Dockerfile which runs a server in Go that will provide API endpoints to perform various tasks with Squid proxy server

### How to Run 
- docker build -t squid-proxy-balancer .
-  docker run -d -p 8080:8080 -p 3128:3128 --name balancer squid-proxy-balancer