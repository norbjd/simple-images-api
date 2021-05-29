run:
	sudo docker-compose up

openapi-ui:
	sudo docker run --rm -it -p 8081:8080 -e SWAGGER_JSON=/openapi.yaml -v `pwd`/api/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui:v3.49.0
