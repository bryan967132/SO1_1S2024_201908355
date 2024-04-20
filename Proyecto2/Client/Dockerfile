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
RUN go build Client.go

# Expone el puerto 3000
EXPOSE 3000

# Comando para ejecutar la aplicación
CMD ["./Client"]

# docker build -t dannyt8355/grpc_client .
# docker run -p 3000:3000 dannyt8355/grpc_client