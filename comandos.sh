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
# VERIFICAR CONTENEDORES EJECUTÁNDOSE
docker ps
# DETENER EJECUCIÓN DE CONTENEDOR
docker stop <ID_CONTENEDOR>

# HT1
# AGREGAR EL MÓDULO
sudo insmod ram.ko
# ELIMINAR MÓDULO
sudo rmmod ram.ko