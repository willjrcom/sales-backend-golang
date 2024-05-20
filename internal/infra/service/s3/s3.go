package s3service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

var (
	bucketName = "sales-backend-golang"
	region     = "sa-east-1"
)

type S3Client struct {
	Client     *s3.Client
	BucketName string
	Region     string
}

func NewS3Client() *S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		panic(err)
	}

	return &S3Client{
		Client:     s3.NewFromConfig(cfg),
		BucketName: bucketName,
		Region:     region,
	}
}

func (s *S3Client) newObject(key string, file *multipart.File) *s3.PutObjectInput {
	return &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(key),
		Body:   *file,
		ACL:    types.ObjectCannedACLPublicRead, // Permite que o objeto seja lido publicamente
	}
}

func (s *S3Client) UploadFile(file *multipart.File) (*string, error) {
	key := uuid.NewString()

	// Configure os parâmetros para o upload
	uploadInput := s.newObject(key, file)

	// Execute o upload
	output, err := s.Client.PutObject(context.TODO(), uploadInput)
	if err != nil {
		return nil, err
	}

	println(output)

	return &key, nil
}

func (s *S3Client) ListObjects() error {
	// Liste os objetos no bucket
	output, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(s.BucketName),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Objects in bucket:")
	for _, obj := range output.Contents {
		fmt.Println(*obj.Key)
	}

	return nil
}

func (s *S3Client) DeleteObject(key string) error {
	objectIds := []types.ObjectIdentifier{
		{Key: aws.String(key)},
	}

	output, err := s.Client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &types.Delete{Objects: objectIds},
	})

	if err != nil {
		log.Printf("Couldn't delete objects from bucket %v. Here's why: %v\n", bucketName, err)
		return err
	}

	log.Printf("Deleted %v objects.\n", len(output.Deleted))
	return nil
}
