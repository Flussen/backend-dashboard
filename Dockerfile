# Usa una imagen de Alpine como base
FROM golang:1.22-alpine

# Instala el gestor de paquetes apk
RUN apk update && apk add --no-cache bash

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de go.mod y go.sum y descarga las dependencias
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copia el resto de los archivos del proyecto al directorio de trabajo
COPY . .

# Copia el archivo .env al directorio de trabajo
COPY .env .

# Compila la aplicación Go
RUN go build -o /myapp

# Puerto en el que se expone la aplicación
EXPOSE 8080

# Ejecuta la aplicación compilada
CMD ["/myapp"]
