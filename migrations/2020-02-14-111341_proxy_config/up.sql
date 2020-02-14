-- Your SQL goes here
CREATE TABLE proxy_config (
    id char(36) DEFAULT (lower(hex(randomblob(4))) || '-' || lower(hex(randomblob(2))) || '-4' || substr(lower(hex(randomblob(2))),2) || '-' || substr('89ab',abs(random()) % 4 + 1, 1) || substr(lower(hex(randomblob(2))),2) || '-' || lower(hex(randomblob(6)))) NOT NULL PRIMARY KEY,
    peer TEXT not null,
    server TEXT not null,
    state INTEGER not null,
    ts TIMESTAMP default current_timestamp not null,
    ts_mod TIMESTAMP default current_timestamp not null
)