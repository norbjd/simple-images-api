DOCKER_CMD=docker # add sudo if you are not in docker group
DOCKER_COMPOSE_CMD=docker-compose # add sudo if you are not in docker group

build:
	$(DOCKER_CMD) build -t "norbjd/simple-images-api:`git rev-parse --short HEAD`" .

push:
	@echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin
	$(DOCKER_CMD) push "norbjd/simple-images-api:`git rev-parse --short HEAD`"
	$(DOCKER_CMD) logout

run:
	$(DOCKER_COMPOSE_CMD) up --build

run-tests: test-setup test test-teardown

test-setup:
	-$(DOCKER_CMD) rm -f 'minio-test'
	$(DOCKER_CMD) run -d --name 'minio-test' -p 9001:9000 \
		-e MINIO_ACCESS_KEY=accessKey -e MINIO_SECRET_KEY=secretKey \
		--entrypoint /bin/sh \
		minio/minio:RELEASE.2021-05-27T22-06-31Z \
		-c "mkdir -p /data/images; minio server /data"
	until curl -f http://localhost:9001/minio/health/live; do sleep 1; done

test-teardown:
	$(DOCKER_CMD) stop -t 0 'minio-test'

test:
	go test -cover

openapi-ui-open:
	# disable-web-security is used to avoid CORS errors
	chromium --incognito --disable-web-security --user-data-dir=`mktemp -d` http://localhost:8081