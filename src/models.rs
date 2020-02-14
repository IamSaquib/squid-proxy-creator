#![allow(proc_macro_derive_resolution_fallback)]

use chrono::NaiveDateTime;
#[derive(Queryable)]
pub struct Proxy {
    pub id: String,
    pub peer: String,
    pub server: String,
    pub state: i32,
    pub ts: NaiveDateTime,
    pub ts_mod: NaiveDateTime,
}
