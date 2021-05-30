run:
	sudo docker-compose up --build

test-setup:
	sudo docker rm -f 'minio-test'
	sudo docker run -d --name 'minio-test' -p 9001:9000 \
		-e MINIO_ACCESS_KEY=accessKey -e MINIO_SECRET_KEY=secretKey \
		--entrypoint /bin/sh \
		minio/minio:RELEASE.2021-05-27T22-06-31Z \
		-c "mkdir -p /data/images; minio server /data"
	until curl -f http://localhost:9001/minio/health/live; do sleep 1; done

test-teardown:
	sudo docker stop -t 0 'minio-test'

test:
	go test -cover

run-tests: test-setup test test-teardown

openapi-ui:
	sudo docker run --rm -it -p 8081:8080 -e SWAGGER_JSON=/openapi.yaml -v `pwd`/api/openapi.yaml:/openapi.yaml swaggerapi/swagger-ui:v3.49.0
