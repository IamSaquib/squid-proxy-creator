# Squid Proxy Creator

A Squid Proxy Dockerfile which runs a server in Go that will provide API endpoints to perform various tasks with Squid proxy server

### How to Run 
- docker build -t squid-proxy-balancer .
- docker run -d --restart=always -p 4128:3128 --volume /srv/docker/squid/cache:/var/spool/squid --name balancer -p 1406:1406 squid-proxy-balancer

### To Create Postgres DB

```sql
    CREATE TABLE "proxy_config" (
        "id" CHAR(36) DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))) NOT NULL PRIMARY KEY, 
        "peers" JSON,
        "host"  TEXT,
        "port" CHAR(5),
        "state" INTEGER,
        "ts" DATETIME DEFAULT current_timestamp,
        "ts_mod" DATETIME DEFAULT current_timestamp
    );
```

```sql
    CREATE TABLE "proxy_port" (
        "port_number" INTEGER,
        "availability" BOOLEAN default 1
    );
insert into proxy_port(port_number) with RECURSIVE n(i) as (SELECT 1026 union all SELECT i +1 from n where i < 65525) SELECT i from n;
```