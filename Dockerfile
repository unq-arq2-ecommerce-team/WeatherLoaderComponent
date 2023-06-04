# Establecer la imagen base de Golang
FROM golang:1.17

# Definir directorio de trabajo
WORKDIR /app

# Copiar archivos de configuración de módulos
COPY go.mod .
COPY go.sum .

# Descargar las dependencias del proyecto
RUN go mod download

# Copiar el código fuente del proyecto y el archivo .env
COPY . .

# Compilar el proyecto teniendo en cuenta la ubicación de main.go en la carpeta 'cmd'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd

# Exponer el puerto que utiliza la aplicación
#EXPOSE 8083

# Ejecutar la aplicación
CMD ["./main"]