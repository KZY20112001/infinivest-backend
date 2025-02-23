package services

import (
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type GenAIService interface {
	GenerateRoboAdvisorPortfolio(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error)
	GenerateAssetAllocations(req dto.AssetAllocationRequest) (dto.AssetAllocationResponse, error)
	GetLatestAssetPrice(symbol string) (float64, error)
}

type genAIServiceImpl struct {
	repo repositories.GenAIRepository
}

type CategoryResult struct {
	Category string
	Assets   dto.Assets
	Error    error
}

func NewGenAIService(gr repositories.GenAIRepository) *genAIServiceImpl {
	return &genAIServiceImpl{repo: gr}
}

func (s *genAIServiceImpl) GenerateRoboAdvisorPortfolio(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error) {
	return s.repo.GeneratePortfolioRecommendation(bankStatement, bankName, toleranceLevel)
}

func (s *genAIServiceImpl) GenerateAssetAllocations(req dto.AssetAllocationRequest) (dto.AssetAllocationResponse, error) {

	resChan := make(chan CategoryResult, len(req.Portfolio))

	var wg sync.WaitGroup
	for category, percentage := range req.Portfolio {
		if category == "cash" {
			continue
		}
		wg.Add(1)
		go func(category string, percentage float64) {
			defer wg.Done()
			if percentage == 0 {
				resChan <- CategoryResult{Category: category, Assets: dto.Assets{Assets: []dto.Asset{}}}
				return
			}
			assets, err := s.repo.GenerateAssetAllocation(category, percentage)
			if err != nil {
				resChan <- CategoryResult{Category: category, Error: err}
			} else {
				resChan <- CategoryResult{Category: category, Assets: assets}
			}

		}(category, percentage)
	}

	wg.Wait()
	close(resChan)
	allocations := make(map[string]dto.Assets)
	var errorList []string

	for res := range resChan {
		if res.Error != nil {
			errorList = append(errorList, fmt.Sprintf("%s: %v", res.Category, res.Error))
		} else {
			allocations[res.Category] = res.Assets
		}
	}
	if len(errorList) > 0 {
		return dto.AssetAllocationResponse{}, fmt.Errorf("failed to generate asset allocation for categories: %v", errorList)
	}
	return dto.AssetAllocationResponse{Allocations: allocations}, nil
}

func (s *genAIServiceImpl) GetLatestAssetPrice(symbol string) (float64, error) {
	return s.repo.GetLatestAssetPrice(symbol)
}
