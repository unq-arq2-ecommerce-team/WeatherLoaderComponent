# WeatherLoaderComponent


## Tecnologías:

- [Golang](https://go.dev/)
- [Gin (WEB API)](https://gin-gonic.com/)
- [MongoDB](https://www.mongodb.com/)

## Prerequisitos:

- Go 1.20 or up / Docker

## Swagger

Instalar swag localmente (se necesita go 1.20 or up)

```
go install github.com/swaggo/swag/cmd/swag@v1.8.10
```

Para actualizar la api doc de swagger, ejecutar en el folder root del repo:

```
swag init -g internal/infrastructure/app.go
```

Luego de levantar la api e ir al endpoint:

```
http://localhost:<port>/docs/index.html
```


## Inicialización y ejecución del proyecto (docker)

### Pasos:

1) Ir a la carpeta root del repositorio

2) Construir el Dockerfile (imagen) del servicio

    ```
    docker build -t weather_loader_component .
    ```

3) Ejecutar la imagen construida.


Tambien, si se desea se puede cambiar las envs por otras de las que estan. Se recomienda utilizar el mismo puerto externo e interno para que funcione correctamente swagger.

```
docker run -p <port>:8083 --env-file ./.env --name weather_loader_component weather_loader_component
```

Nota: agregar "-d" si se quiere ejecutar como deamon

Ejemplo:

```
docker run -d -p 8083:8083 --env-file ./.env --name weather_loader_component weather_loader_component
```

4) En un browser, abrir swagger del servicio en el siguiente url:

`http://localhost:<port>/docs/index.html`

5) Probar el endpoint health check y debe retornar ok

6) La API esta disponible para ser utilizada