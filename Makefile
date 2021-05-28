build:
	sudo docker build -t 'norbjd/simple-images-api:dev' .

run:
	sudo docker run --rm -it --name 'simple-images-api' -p 8080:8080 'norbjd/simple-images-api:dev'
