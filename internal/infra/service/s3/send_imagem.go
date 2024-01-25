package s3

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

func UploadToS3(file *multipart.File) (*string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil, err
	}

	// Crie um cliente S3
	client := s3.NewFromConfig(cfg)
	key := uuid.NewString()

	// Configure os parâmetros para o upload
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String("sales-backend-golang"),
		Key:    aws.String(key),
		Body:   *file,
		ACL:    types.ObjectCannedACLPublicRead, // Permite que o objeto seja lido publicamente
	}

	// Execute o upload
	output, err := client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		return nil, err
	}

	println(output)
	fmt.Printf("Upload bem-sucedido. A imagem está disponível em: https://%s.s3.amazonaws.com/%s\n")

	return &key, nil
}
