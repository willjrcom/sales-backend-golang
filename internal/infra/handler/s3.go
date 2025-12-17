package handlerimpl

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/willjrcom/sales-backend-go/bootstrap/database"
	"github.com/willjrcom/sales-backend-go/bootstrap/handler"
	s3dto "github.com/willjrcom/sales-backend-go/internal/infra/dto/s3"
	jsonpkg "github.com/willjrcom/sales-backend-go/pkg/json"
)

func NewHandlerS3() *handler.Handler {
	c := chi.NewRouter()

	route := "/s3"

	c.With().Group(func(c chi.Router) {
		c.Post("/presign", handlerGeneratePresignedURL)
		c.Options("/presign", handlerGeneratePresignedURL)
		c.Post("/upload", handlerUploadImage)
	})

	unprotectedRoutes := []string{}
	return handler.NewHandler(route, c, unprotectedRoutes...)
}

func handlerGeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	// Adicionar headers CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, access-token")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := r.Context()

	// Parse request body
	var req s3dto.PresignRequestDTO
	if err := jsonpkg.ParseBody(r, &req); err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusBadRequest, err)
		return
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, fmt.Errorf("failed to load AWS config: %w", err))
		return
	}

	// Create S3 client
	s3Client := awss3.NewFromConfig(cfg)

	// Generate unique key with schema_name in the path
	var key string
	if req.SchemaName != "" {
		key = fmt.Sprintf("images/%s/%d-%s", req.SchemaName, time.Now().Unix(), req.Filename)
	} else {
		key = fmt.Sprintf("images/public/%d-%s", time.Now().Unix(), req.Filename)
	}

	// Create presigned URL
	presignClient := awss3.NewPresignClient(s3Client)
	presignResult, err := presignClient.PresignPutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:         aws.String(key),
		ContentType: aws.String(req.ContentType),
		ACL:         "public-read",
	}, awss3.WithPresignExpires(time.Minute)) // URL válida por 1 minuto

	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, fmt.Errorf("failed to generate presigned URL: %w", err))
		return
	}

	// Generate public URL
	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET_NAME")
	publicUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)

	// Return response
	response := s3dto.PresignResponseDTO{
		URL:       presignResult.URL,
		Key:       key,
		PublicUrl: publicUrl,
	}

	jsonpkg.ResponseJson(w, r, http.StatusOK, response)
}

func handlerUploadImage(w http.ResponseWriter, r *http.Request) {
	// Limite de 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, fmt.Errorf("erro ao processar upload: %w", err))
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, fmt.Errorf("arquivo não encontrado: %w", err))
		return
	}
	defer file.Close()

	schemaName, _ := database.GetCurrentSchema(r.Context())

	// Gera nome único
	filename := handler.Filename
	sanitizedName := filename
	key := ""
	if schemaName != "" {
		key = fmt.Sprintf("images/%s/%d-%s", schemaName, time.Now().Unix(), sanitizedName)
	} else {
		key = fmt.Sprintf("images/public/%d-%s", time.Now().Unix(), sanitizedName)
	}

	// Carrega config AWS
	ctx := r.Context()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, fmt.Errorf("erro AWS config: %w", err))
		return
	}
	s3Client := awss3.NewFromConfig(cfg)

	bucket := os.Getenv("S3_BUCKET_NAME")

	// Upload para o S3
	_, err = s3Client.PutObject(ctx, &awss3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(handler.Header.Get("Content-Type")),
		ACL:         "public-read",
	})
	if err != nil {
		jsonpkg.ResponseErrorJson(w, r, http.StatusInternalServerError, fmt.Errorf("erro ao enviar para o S3: %w", err))
		return
	}

	region := os.Getenv("AWS_REGION")
	publicUrl := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, key)

	response := s3dto.PresignResponseDTO{
		Key:       key,
		PublicUrl: publicUrl,
	}
	jsonpkg.ResponseJson(w, r, http.StatusCreated, response)
}
