package minio

// Estrutura para request de deleção múltipla
type DeleteImagesRequest struct {
	ImageURLs []string `json:"image_urls"`
}

// Estrutura para resultado de deleção
type DeleteImageResult struct {
	ImageURL string `json:"image_url"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
	Key      string `json:"key,omitempty"`
}

// Job para deleção de imagem
type DeleteImageJob struct {
	ImageURL string
	Index    int
}

// Resultado do job de deleção
type DeleteImageJobResult struct {
	Result DeleteImageResult
	Index  int
}
