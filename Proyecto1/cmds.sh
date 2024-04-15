# AGREGAR EL MÓDULO
sudo insmod m.ko
# ELIMINAR MÓDULO
sudo rmmod m.ko

# DOCKER
# CREAR IMAGEN DE BACK
docker build -t back_t1 .
# CREAR IMAGEN DE FRONT
docker build -t front_t1 .
# VERIFICAR IMÁGENES EXISTENTES
docker images
# EJECUTAR IMAGEN BACK
docker run -d -p 8000:8000 back_t1
# EJECUTAR IMAGEN FRONT
docker run -d -p 3000:3000 front_t1
# ELIMINAR CONTENEDOR
docker rm <ID_CONTENEDOR>
# ELIMINAR IMAGEN
docker rmi <ID_IMAGEN>
# VERIFICAR CONTENEDORES EJECUTÁNDOSE
docker ps
# DETENER EJECUCIÓN DE CONTENEDOR
docker stop <ID_CONTENEDOR>
# DOCKER-COMPOSE
docker-compose up -d

# BASE DE DATOS
# CONECTARSE A LA BD CON DOCKER-COMPOSE
sudo docker-compose exec mysql mysql -uroot -p
# CONECTARSE A LA BD
sudo mysql -u root -p
# DETENER MYSQL-SERVER LOCAL
sudo service mysql stop