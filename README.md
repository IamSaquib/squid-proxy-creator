# Squid Proxy Creator

A Squid Proxy Dockerfile which runs a server in Go that will provide API endpoints to perform various tasks with Squid proxy server

### How to Run 
- docker build -t squid-proxy-balancer .
- docker run -d --restart=always -p 4128:3128 --volume /srv/docker/squid/cache:/var/spool/squid --name balancer -p 1406:1406 squid-proxy-balancer

### To Create Postgres DB

```sql
    create table "proxy_config" (
        "id" char(36) DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))) NOT NULL PRIMARY KEY, 
        "peers" json_array,
        "server" TEXT,
        "state" INTEGER,
        "ts" DATETIME default current_timestamp,
        "ts_mod" DATETIME default current_timestamp
    );
```

