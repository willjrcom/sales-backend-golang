package s3

type PresignRequestDTO struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	SchemaName  string `json:"schema_name"`
}
