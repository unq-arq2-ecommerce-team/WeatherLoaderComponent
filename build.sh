#!/bin/bash

# Variables para controlar la ejecución de los tests y el contenedor
run_tests=false
run_container=false

# Iterar sobre los argumentos y establecer las variables según las opciones proporcionadas
for arg in "$@"
do
  if [ "$arg" == "--test" ]; then
    run_tests=true
  elif [ "$arg" == "--run" ]; then
    run_container=true
  fi
done

# Si se proporcionó la opción '--test', ejecutar los tests antes de construir la imagen
if [ "$run_tests" == "true" ]; then
  echo "Ejecutando los tests..."
  ./run-tests.sh
  TEST_EXIT_CODE=$?
  if [ $TEST_EXIT_CODE -ne 0 ]; then
    echo "Algunos tests fallaron. Abortando la construcción de la imagen Docker."
    exit 1
  fi
fi

# Construir la imagen de Docker
docker build -t weather-service .

# Si se proporcionó la opción '--run', ejecutar el contenedor después de construir la imagen
if [ "$run_container" == "true" ]; then
  echo "Ejecutando el contenedor..."
  docker run -p 8080:8080 weather-service
fi