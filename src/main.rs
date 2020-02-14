use std::net::SocketAddr;
use hyper::{Body, Request, Response, Server, Method, StatusCode};
use hyper::service::{make_service_fn, service_fn};

mod endpoints;
use endpoints::get_proxy;

#[macro_use]
extern crate diesel;
extern crate dotenv;

pub mod models;
pub mod schema;

use diesel::prelude::*;
use dotenv::dotenv;
use std::env;
use self::models::*;

pub fn establish_connection() -> SqliteConnection {
    dotenv().ok();

    let database_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    SqliteConnection::establish(&database_url)
        .unwrap_or_else(|_| panic!("Error connecting to {}", database_url))
}

pub fn show_proxy() {
    use schema::proxy_config::dsl::*;
    let connection = establish_connection();
    let results = proxy_config
        .load::<Proxy>(&connection)
        .expect("Error loadig proxy");
        println!("There there!");
    for pr in results {
        println!("{:?}", pr.id);
    }
}
async fn handle(req: Request<Body>) -> Result<Response<Body>, hyper::Error> {
    match (req.method(), req.uri().path()) {
        // Health to check up the server status
        (&Method::GET, "/health")=> Ok(Response::new(Body::from(
            "Server Working",
        ))),

        // Proxy to get content of a particular proxy
        (&Method::GET, "/proxy")=> get_proxy(),

        // Return the 404 Not Found for other routes.
        _ => {
            let mut not_found = Response::default();
            *not_found.status_mut() = StatusCode::NOT_FOUND;
            Ok(not_found)
        }
    }
}

async fn shutdown_signal() {
    // Wait for the CTRL+C signal
    tokio::signal::ctrl_c()
        .await
        .expect("failed to install CTRL+C signal handler");
}


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error + Send + Sync>> {
    let addr = SocketAddr::from(([127, 0, 0, 1], 8080));

    show_proxy();
    let make_svc = make_service_fn(move |_conn| async {
        Ok::<_, hyper::Error>(service_fn(handle))
    });

    let server = Server::bind(&addr).serve(make_svc);

    let graceful = server.with_graceful_shutdown(shutdown_signal());
    if let Err(e) = graceful.await {
        eprintln!("server error: {}", e);
    }
    Ok(())
}