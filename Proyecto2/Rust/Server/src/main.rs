use rocket::response::status::BadRequest;
use rocket::serde::json::{json, Value as JsonValue};
use rocket::serde::json::Json;
use rocket::config::SecretKey;
use rocket_cors::{AllowedOrigins, CorsOptions};

#[derive(rocket::serde::Deserialize)]
struct Data {
    name: String,
    album: String,
    year: String,
    rank: String,
}

#[rocket::get("/")]
fn get_data() -> Result<String, BadRequest<String>> {
    Ok(JsonValue::from(json!({
        "message": format!("Server is running!")
    })).to_string())
}

#[rocket::post("/data", data = "<data>")]
fn receive_data(data: Json<Data>) -> Result<String, BadRequest<String>> {
    let received_data = data.into_inner();
    let response = JsonValue::from(json!({
        "message": format!("Received data: name: {}, album: {}, year: {}, rank: {}", received_data.name, received_data.album, received_data.year, received_data.rank)
    }));
    Ok(response.to_string())
}

#[rocket::main]
async fn main() {
    let secret_key = SecretKey::generate(); // Genera una nueva clave secreta

    // Configuración de opciones CORS
    let cors = CorsOptions::default()
        .allowed_origins(AllowedOrigins::all())
        .to_cors()
        .expect("failed to create CORS fairing");

    let config = rocket::Config {
        address: "0.0.0.0".parse().unwrap(),
        port: 4001,
        secret_key: secret_key.unwrap(), // Desempaqueta la clave secreta generada
        ..rocket::Config::default()
    };

    // Montar la aplicación Rocket con el middleware CORS
    rocket::custom(config)
        .attach(cors)
        .mount("/", rocket::routes![get_data, receive_data])
        .launch()
        .await
        .unwrap();
}