docker rmi dannyt8355/rust-client:latest
docker rmi dannyt8355/rust-server:latest
cd Rust/Client
docker build -t dannyt8355/rust-client:latest .
cd ../Server
docker build -t dannyt8355/rust-server:latest .
cd ../
docker push dannyt8355/rust-client:latest
docker push dannyt8355/rust-server:latest