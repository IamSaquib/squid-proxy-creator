table! {
    proxy_config (id) {
        id -> Text,
        peer -> Text,
        server -> Text,
        state -> Integer,
        ts -> Timestamp,
        ts_mod -> Timestamp,
    }
}
