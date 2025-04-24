package bootstrap

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func SetUpMinio() *minio.Client{
	minioClient, err := minio.New(MINIO_HOST, &minio.Options{
		Creds:  credentials.NewStaticV4(MINIO_USER, MINIO_PASSWORD, ""),
		Secure: true,
	})
	if err != nil{
		log.Fatalln(err)
	}
	return minioClient
}