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
	GeneratePortfolioRecommendation(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error)
	GenerateAssetAllocation(category string, percentage float64) (dto.Assets, error)
	GetLatestAssetPrice(symbol string) (float64, error)
}

type flaskMicroservice struct {
	client  *http.Client
	baseURL string
}

func NewFlaskMicroservice(baseURL string) *flaskMicroservice {
	return &flaskMicroservice{
		client:  &http.Client{Timeout: 60 * time.Second},
		baseURL: baseURL,
	}
}

func (r *flaskMicroservice) GeneratePortfolioRecommendation(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error) {
	file, err := bankStatement.Open()
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to open bank statement file: %w", err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fileWriter, err := writer.CreateFormFile("bank_statement", bankStatement.Filename)
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to copy file contents: %w", err)
	}

	err = writer.WriteField("bank_name", bankName)
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to add bank name: %w", err)
	}

	err = writer.WriteField("risk_tolerance_level", toleranceLevel)
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to add risk tolerance level: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := fmt.Sprintf("%s/robo-advisor/generate/categories", r.baseURL)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := r.client.Do(req)
	if err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("microservice error: %s", errorResponse["error"])
		}
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("microservice error: status code %d", resp.StatusCode)
	}

	var roboAdvisorRecommendation dto.RoboAdvisorRecommendationResponse
	if err = json.NewDecoder(resp.Body).Decode(&roboAdvisorRecommendation); err != nil {
		return dto.RoboAdvisorRecommendationResponse{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return roboAdvisorRecommendation, nil
}

func (r *flaskMicroservice) GenerateAssetAllocation(category string, percentage float64) (dto.Assets, error) {
	url := fmt.Sprintf("%s/robo-advisor/generate/assets", r.baseURL)
	body, err := json.Marshal(map[string]interface{}{
		"category":   category,
		"percentage": percentage,
	})
	if err != nil {
		return dto.Assets{}, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return dto.Assets{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return dto.Assets{}, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var assets dto.Assets
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return dto.Assets{}, err
	}
	return assets, nil
}

func (r *flaskMicroservice) GetLatestAssetPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("%s/assets/latest-price/%s", r.baseURL, symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var price float64
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return 0, err
	}
	return price, nil
}
