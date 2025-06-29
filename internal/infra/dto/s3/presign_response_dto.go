package s3

type PresignResponseDTO struct {
	URL       string `json:"url"`
	Key       string `json:"key"`
	PublicUrl string `json:"public_url"`
}
