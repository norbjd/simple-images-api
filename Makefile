build:
	sudo docker build -t 'norbjd/simple-images-api:dev' .

run:
	sudo docker run --rm -it --name 'simple-images-api' -p 8080:8080 'norbjd/simple-images-api:dev'

openapi-ui:
	sudo docker run --rm -it -p 8081:8080 -e SWAGGER_JSON=/openapi.yaml -v `pwd`/api/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui:v3.49.0
