use hyper::{Response, Body};

pub fn get_proxy() -> Result<Response<Body>, hyper::Error> {
    
    
    Ok(Response::new(Body::from(
            "Proxy Working",
        )))
}