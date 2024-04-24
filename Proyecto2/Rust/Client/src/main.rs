use rocket::{serde::json::Json};
use rocket::post;
use reqwest::Client;
use serde::{Deserialize, Serialize};
use rocket::config::SecretKey;
use rocket_cors::{AllowedOrigins, CorsOptions};

#[derive(Debug, Serialize, Deserialize)]
struct Data {
    name: String,
    album: String,
    year: String,
    rank: String,
}

#[post("/insert", data = "<data>")]
async fn insert(data: Json<Data>) -> String {
    let client = Client::new();
    let server_url = "http://localhost:4001/data";
    let response = client.post(server_url).json(&data.into_inner()).send().await;

    match response {
        Ok(_) => "Data sent successfully!".to_string(),
        Err(e) => format!("Failed to send data: {}", e),
    }
}

#[rocket::main]
async fn main() {
    // --------------------------------------------------------------------------
    let secret_key = SecretKey::generate(); // Genera una nueva clave secreta

    // Configuración de opciones CORS
    let cors = CorsOptions::default()
        .allowed_origins(AllowedOrigins::all())
        .to_cors()
        .expect("failed to create CORS fairing");

    let config = rocket::Config {
        address: "0.0.0.0".parse().unwrap(),
        port: 4000,
        secret_key: secret_key.unwrap(), // Desempaqueta la clave secreta generada
        ..rocket::Config::default()
    };
    //----------------------------------------------------------------------------
    rocket::custom(config)
    .attach(cors)
    .mount("/rust", rocket::routes![insert])
    .launch()
    .await
    .unwrap();
}