# Usa una imagen base de Go
FROM golang:1.22.2-alpine

# Establece el directorio de trabajo en /app
WORKDIR /app

# Copia los archivos necesarios para construir la aplicación
COPY go.mod go.sum ./
RUN go mod download

# Copia el resto de la aplicación al directorio de trabajo
COPY . .

# Compila la aplicación
RUN go build Server.go

# Expone el puerto 3001
EXPOSE 3001

# Comando para ejecutar la aplicación
CMD ["./Server"]

# docker build -t proyecto2_grpc_server .
# docker run -p 3001:3001 dannyt8355/grpc_server