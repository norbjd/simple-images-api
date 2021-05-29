package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func getMinioRepository() MinioRepository {
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	minioAccessKey := os.Getenv("MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("MINIO_SECRET_KEY")
	minioBucketName := os.Getenv("MINIO_BUCKET")

	if minioEndpoint == "" || minioAccessKey == "" || minioSecretKey == "" || minioBucketName == "" {
		log.Fatal("Missing one or more env vars : MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY, MINIO_BUCKET")
	}

	minioClient, err := NewMinioClient(minioEndpoint, minioAccessKey, minioSecretKey)
	if err != nil {
		log.Fatal("Cannot init minio client : ", err)
	}

	bucketExists, err := minioClient.BucketExists(minioBucketName)
	if err != nil {
		log.Fatal("Error while checking if " + minioBucketName + " bucket exists : ", err)
	}
	if !bucketExists {
		log.Fatal("Bucket \"" + minioBucketName + "\" does not exist : create it before running app")
	}

	return MinioRepository{bucketName: minioBucketName, client: minioClient}
}

func main() {
	App{repository: getMinioRepository()}.run()
}
