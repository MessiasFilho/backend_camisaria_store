package minio

import (
	"backend_camisaria_store/config"
	"backend_camisaria_store/schemas"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

// Helper function to convert JSON []byte to []string
func JsonToStringSlice(jsonData []byte) []string {
	if jsonData == nil || len(jsonData) == 0 {
		return []string{}
	}

	var result []string
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return []string{}
	}

	return result
}

// Estrutura para resultado do upload de imagem
type ImageUploadResult struct {
	FileHeader *multipart.FileHeader
	PublicURL  string
	ObjectName string
	Error      error
	Order      int
}

// Worker para processamento paralelo de uploads
func processImageUpload(jobChan <-chan ImageUploadJob, resultChan chan<- ImageUploadResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobChan {
		// Validar arquivo
		if err := ValidateImageFile(job.FileHeader); err != nil {
			resultChan <- ImageUploadResult{
				FileHeader: job.FileHeader,
				Error:      err,
				Order:      job.Order,
			}
			continue
		}

		// Fazer upload para MinIO
		publicURL, objectName, err := UploadProductImage(job.FileHeader, job.ProductName, "products", job.ProductID)

		resultChan <- ImageUploadResult{
			FileHeader: job.FileHeader,
			PublicURL:  publicURL,
			ObjectName: objectName,
			Error:      err,
			Order:      job.Order,
		}
	}
}

// Estrutura para job de upload
type ImageUploadJob struct {
	FileHeader  *multipart.FileHeader
	ProductName string
	ProductID   uint64
	Order       int
}

func UploadImgesProduct(c *fiber.Ctx) error {

	// 1. Validar se o ID foi fornecido
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "ID do produto é obrigatório",
			"message": "O parâmetro 'id' não foi fornecido na URL",
		})
	}

	// 2. Validar se o ID é um número válido e positivo
	productID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":    "ID do produto inválido",
			"message":  "O ID deve ser um número inteiro positivo",
			"received": idParam,
		})
	}

	if productID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "ID do produto inválido",
			"message": "O ID do produto deve ser maior que zero",
		})
	}

	produto := schemas.Products{}
	result := config.DB.Where("id = ?", productID).First(&produto)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":      "Produto não encontrado",
				"message":    "Não existe produto com o ID informado",
				"product_id": productID,
			})
		}
		// Erro de banco de dados
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Erro interno do servidor",
			"message": "Erro ao consultar produto no banco de dados",
		})
	}

	if !produto.IsActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      "Produto inativo",
			"message":    "Não é possível fazer upload de imagens para produtos inativos",
			"product_id": productID,
		})
	}

	// ====================
	// Validação: Verificar limite máximo de imagens (5)
	// ====================

	// Contar imagens já existentes
	existingImages := JsonToStringSlice(produto.Images)
	currentImageCount := len(existingImages)

	const maxImagesPerProduct = 5
	if currentImageCount >= maxImagesPerProduct {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Limite de imagens excedido",
			"message": "Este produto já possui o número máximo de imagens permitidas",
		})
	}

	// Obter arquivos multipart do formulário
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro no formulário",
			"message": "Erro ao processar dados do formulário multipart",
		})
	}

	files := form.File["images"] // Campo esperado: "images"
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Nenhum arquivo enviado",
			"message": "Envie pelo menos uma imagem no campo 'images'",
		})
	}

	// Limite de arquivos por upload (considerando imagens já existentes)
	maxFilesAllowed := maxImagesPerProduct - currentImageCount
	if len(files) > maxFilesAllowed {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":              "Muitos arquivos",
			"message":            "Máximo de " + strconv.Itoa(maxFilesAllowed) + " imagens permitidas neste upload",
			"current_images":     currentImageCount,
			"max_allowed_upload": maxFilesAllowed,
			"max_total_images":   maxImagesPerProduct,
		})
	}

	// ====================
	// Processar imagens em paralelo (alta performance)
	// ====================

	// Preparar imagens existentes
	existingImages = JsonToStringSlice(produto.Images)
	newImageURLs := make([]string, len(existingImages))
	copy(newImageURLs, existingImages)

	// Configurar concorrência baseada no número de CPUs disponíveis
	numWorkers := runtime.NumCPU()
	if len(files) < numWorkers {
		numWorkers = len(files) // Não criar mais workers que arquivos
	}
	if numWorkers > 4 {
		numWorkers = 4 // Limitar a no máximo 4 workers simultâneos
	}

	// Canais para comunicação entre goroutines
	jobChan := make(chan ImageUploadJob, len(files))
	resultChan := make(chan ImageUploadResult, len(files))

	// Iniciar workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processImageUpload(jobChan, resultChan, &wg)
	}

	// Enviar jobs para processamento
	go func() {
		defer close(jobChan)
		for i, fileHeader := range files {
			jobChan <- ImageUploadJob{
				FileHeader:  fileHeader,
				ProductName: produto.Name,
				ProductID:   productID,
				Order:       len(existingImages) + i,
			}
		}
	}()

	// Fechar resultChan quando todos os workers terminarem
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados
	var uploadedImages []fiber.Map
	var errors []fiber.Map
	resultsCount := 0

	for result := range resultChan {
		resultsCount++

		if result.Error != nil {
			errors = append(errors, fiber.Map{
				"file":  result.FileHeader.Filename,
				"error": result.Error.Error(),
			})
			continue
		}

		// Upload bem-sucedido
		newImageURLs = append(newImageURLs, result.PublicURL)
		uploadedImages = append(uploadedImages, fiber.Map{
			"url":        result.PublicURL,
			"filename":   result.FileHeader.Filename,
			"size":       result.FileHeader.Size,
			"order":      result.Order,
			"object_key": result.ObjectName,
		})
	}

	// ====================
	// Atualizar produto com novas imagens
	// ====================

	if len(uploadedImages) > 0 {
		// Converter slice de strings para JSON
		imagesJSON, err := json.Marshal(newImageURLs)
		if err != nil {
			errors = append(errors, fiber.Map{
				"error": "Erro ao converter imagens para JSON: " + err.Error(),
			})
		} else {
			// Usar SQL direto para atualizar apenas o campo Images
			result := config.DB.Model(&schemas.Products{}).
				Where("id = ?", productID).
				Update("images", imagesJSON)

			if result.Error != nil {
				// Se falhar ao atualizar, adicionar erro
				errors = append(errors, fiber.Map{
					"error": "Erro ao atualizar produto com novas imagens: " + result.Error.Error(),
				})
			} else if result.RowsAffected == 0 {
				// Se não afetou nenhuma linha, produto não foi encontrado
				errors = append(errors, fiber.Map{
					"error": "Produto não encontrado para atualização",
				})
			}
		}
	}

	// ====================
	// Retornar resposta
	// ====================

	message := "Upload processado"
	response := fiber.Map{
		"message": message,
	}

	if len(errors) > 0 {
		response["errors"] = errors
		message += " com alguns erros"
		response["message"] = message
		return c.Status(fiber.StatusPartialContent).JSON(response)
	}

	response["message"] = "Todas as imagens foram enviadas com sucesso"
	return c.Status(fiber.StatusCreated).JSON(response)
}

func DeleteImagesMinio(c *fiber.Ctx) error {

	// ====================
	// Parse do Request
	// ====================

	var req DeleteImagesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Erro ao processar dados da requisição",
			"message": err.Error(),
		})
	}

	// Validações básicas
	if len(req.ImageURLs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Lista de URLs vazia",
			"message": "Envie pelo menos uma URL de imagem para deletar",
		})
	}

	if len(req.ImageURLs) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Muitas URLs",
			"message": "Máximo de 10 imagens por vez",
		})
	}

	// ====================
	// Processamento Paralelo
	// ====================

	// Configurar concorrência
	numWorkers := runtime.NumCPU()
	if len(req.ImageURLs) < numWorkers {
		numWorkers = len(req.ImageURLs)
	}
	if numWorkers > 4 {
		numWorkers = 4 // Máximo 4 workers simultâneos
	}

	// Canais para comunicação
	jobChan := make(chan DeleteImageJob, len(req.ImageURLs))
	resultChan := make(chan DeleteImageJobResult, len(req.ImageURLs))

	// Iniciar workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go deleteImageWorker(jobChan, resultChan, &wg)
	}

	// Enviar jobs
	go func() {
		defer close(jobChan)
		for i, imageURL := range req.ImageURLs {
			jobChan <- DeleteImageJob{
				ImageURL: imageURL,
				Index:    i,
			}
		}
	}()

	// Fechar resultChan quando workers terminarem
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Coletar resultados
	results := make([]DeleteImageResult, len(req.ImageURLs))
	processedCount := 0

	for jobResult := range resultChan {
		results[jobResult.Index] = jobResult.Result
		processedCount++
	}

	// ====================
	// Preparar Resposta
	// ====================

	successCount := 0
	errorCount := 0

	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			errorCount++
		}
	}

	response := fiber.Map{
		"message":       fmt.Sprintf("Processamento concluído: %d sucesso(s), %d erro(s)", successCount, errorCount),
		"total_images":  len(req.ImageURLs),
		"success_count": successCount,
		"error_count":   errorCount,
		"results":       results,
	}

	if errorCount > 0 {
		return c.Status(fiber.StatusPartialContent).JSON(response)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Worker para processamento paralelo de deleções
func deleteImageWorker(jobChan <-chan DeleteImageJob, resultChan chan<- DeleteImageJobResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobChan {
		result := DeleteImageJobResult{
			Index: job.Index,
			Result: DeleteImageResult{
				ImageURL: job.ImageURL,
				Success:  false,
			},
		}

		// Validações de segurança
		if job.ImageURL == "" {
			result.Result.Error = "URL da imagem é obrigatória"
			resultChan <- result
			continue
		}

		// Validar formato da URL
		if !strings.Contains(job.ImageURL, config.BunkedName) && !strings.Contains(job.ImageURL, "minio") {
			result.Result.Error = "URL inválida: não pertence ao domínio MinIO"
			resultChan <- result
			continue
		}

		// Extrair key da URL
		key, err := ObjectKeyFormUrl(job.ImageURL)
		if err != nil {
			result.Result.Error = fmt.Sprintf("falha ao extrair key: %v", err)
			resultChan <- result
			continue
		}

		// Verificar se a key não está vazia
		if key == "" {
			result.Result.Error = "key do objeto não pôde ser extraída da URL"
			resultChan <- result
			continue
		}

		// Verificar se pertence ao diretório de produtos
		if !strings.Contains(key, "products/") {
			result.Result.Error = "key inválida: não pertence ao diretório de produtos"
			resultChan <- result
			continue
		}

		// Executar deleção com timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		err = config.MinioClient.RemoveObject(
			ctx,
			config.BunkedName,
			key,
			minio.RemoveObjectOptions{},
		)
		cancel() // Cancelar contexto após uso

		if err != nil {
			result.Result.Error = fmt.Sprintf("falha ao deletar: %v", err)
			resultChan <- result
			continue
		}

		// Sucesso
		result.Result.Success = true
		result.Result.Key = key
		resultChan <- result
	}
}
