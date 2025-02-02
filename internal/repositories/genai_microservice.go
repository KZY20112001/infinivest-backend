package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
)

type GenAIRepository interface {
	GetPortfolioRecommendation(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error)
}

type flaskMicroservice struct {
	client  *http.Client
	baseURL string
}

func NewFlaskMicroservice(baseURL string) *flaskMicroservice {
	return &flaskMicroservice{
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: baseURL,
	}
}

func (r *flaskMicroservice) GetPortfolioRecommendation(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error) {
	file, err := bankStatement.Open()
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to open bank statement file: %w", err)
	}
	defer file.Close()

	var dto dto.RoboAdvisorRecommendationResponse

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile("bank_statement", bankStatement.Filename)
	if err != nil {
		return dto, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return dto, fmt.Errorf("failed to copy file contents: %w", err)
	}

	err = writer.WriteField("risk_tolerance_level", toleranceLevel)
	if err != nil {
		return dto, fmt.Errorf("failed to add risk tolerance level: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return dto, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := fmt.Sprintf("%s/robo-advisor/%s/generate", r.baseURL, bankName)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return dto, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := r.client.Do(req)
	if err != nil {
		return dto, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			return dto, fmt.Errorf("microservice error: %s", errorResponse["error"])
		}
		return dto, fmt.Errorf("microservice error: status code %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&dto)
	if err != nil {
		return dto, fmt.Errorf("failed to decode response: %w", err)
	}

	return dto, nil
}
