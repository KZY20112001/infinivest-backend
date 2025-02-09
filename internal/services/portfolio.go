package services

import (
	"fmt"
	"mime/multipart"
	"sync"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type PortfolioService interface {
	GenerateRoboAdvisorPortfolio(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error)
	GenerateAssetAllocations(dto dto.AssetAllocationRequest) (dto.AssetAllocationResponse, error)
}

type portfolioServiceImpl struct {
	repo repositories.GenAIRepository
}

type CategoryResult struct {
	Category string
	Assets   dto.Assets
	Error    error
}

func NewPortfolioServiceImpl(r repositories.GenAIRepository) *portfolioServiceImpl {
	return &portfolioServiceImpl{repo: r}
}

func (s *portfolioServiceImpl) GenerateRoboAdvisorPortfolio(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error) {
	return s.repo.GeneratePortfolioRecommendation(bankStatement, bankName, toleranceLevel)
}

func (s *portfolioServiceImpl) GenerateAssetAllocations(req dto.AssetAllocationRequest) (dto.AssetAllocationResponse, error) {

	resChan := make(chan CategoryResult, len(req.Portfolio))

	var wg sync.WaitGroup
	for category, percentage := range req.Portfolio {

		wg.Add(1)
		go s.generateAssetAllocation(category, percentage, &wg, resChan)
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

func (s *portfolioServiceImpl) generateAssetAllocation(category string, percentage float64, wg *sync.WaitGroup, resChan chan<- CategoryResult) {
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
}
