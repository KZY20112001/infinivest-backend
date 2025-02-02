package services

import (
	"mime/multipart"

	"github.com/KZY20112001/infinivest-backend/internal/dto"
	"github.com/KZY20112001/infinivest-backend/internal/repositories"
)

type PortfolioService interface {
	GenerateRoboAdvisorPortfolio(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error)
}

type portfolioServiceImpl struct {
	repo repositories.GenAIRepository
}

func NewPortfolioServiceImpl(r repositories.GenAIRepository) *portfolioServiceImpl {
	return &portfolioServiceImpl{repo: r}
}

func (s *portfolioServiceImpl) GenerateRoboAdvisorPortfolio(bankStatement *multipart.FileHeader, bankName, toleranceLevel string) (dto.RoboAdvisorRecommendationResponse, error) {
	return s.repo.GetPortfolioRecommendation(bankStatement, bankName, toleranceLevel)
}
