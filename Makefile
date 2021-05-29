build:
	sudo docker build -t 'norbjd/simple-images-api:dev' .

run: build
	sudo docker run --rm -it --name 'simple-images-api' -p 8080:8080 -e LOG_LEVEL=debug -e MINIO_ENDPOINT=localhost:9000 -e MINIO_ACCESS_KEY=M1N104cC355K3y -e MINIO_SECRET_KEY=MInI053cR3tK3y 'norbjd/simple-images-api:dev'

openapi-ui:
	sudo docker run --rm -it -p 8081:8080 -e SWAGGER_JSON=/openapi.yaml -v `pwd`/api/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui:v3.49.0

minio:
	sudo docker run --rm -it --name 'minio' -p 9000:9000 -e MINIO_ACCESS_KEY=M1N104cC355K3y -e MINIO_SECRET_KEY=MInI053cR3tK3y -v `pwd`/minio_data:/data minio/minio:RELEASE.2021-05-27T22-06-31Z server /data