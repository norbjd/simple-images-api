package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func getMinioRepository() MinioRepository {
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")

	if minioEndpoint == "" || minioAccessKey == "" || minioSecretKey == "" {
		log.Fatal("Missing one or more env vars : MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY")
	}

	minioClient, err := NewMinioClient(minioEndpoint, minioAccessKey, minioSecretKey)
	if err != nil {
		log.Fatal("Cannot init minio client", err)
	}

	return MinioRepository{client: minioClient}
}

func main() {
	App{repository: getMinioRepository()}.run()
}
