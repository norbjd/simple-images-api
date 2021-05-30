run:
	sudo docker-compose up --build

run-tests: test-setup test test-teardown

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

openapi-ui-open:
	# disable-web-security is used to avoid CORS errors
	chromium --incognito --disable-web-security --user-data-dir=`mktemp -d` http://localhost:8081