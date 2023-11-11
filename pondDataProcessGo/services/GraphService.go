package services

import (
	"pondDataProcessGo/models"
	"pondDataProcessGo/repositories"
)

type GraphService struct {
	graphRepository *repositories.GraphRepository
}

func NewGraphService() *GraphService {
	return &GraphService{
		graphRepository: repositories.NewGraphRepository(),
	}
}

func (s *GraphService) MergeTwitterAccount(users []models.User) {
	err := s.graphRepository.MergeTwitterAccount(users)
	if err != nil {
		return
	}
}

func (s *GraphService) MergeTwitterFollowing(user models.User, followUsers []models.User, uType bool) {
	err := s.graphRepository.MergeTwitterFollowing(user, followUsers, uType)
	if err != nil {
		return
	}
}

func (s *GraphService) MergeTwitter2Person(newPersons []models.User) []int64 {
	return s.graphRepository.MergeTwitter2Person(newPersons)
}

func (s *GraphService) MergePerson2Person(user models.User, rType bool) {
	s.graphRepository.MergePerson2Person(user, rType)
}

func (s *GraphService) UpdateComputeRelationStrength(user models.User) ([]int64, []int64, []float64) {
	return s.graphRepository.UpdateComputeRelationStrength(user)
}

func (s *GraphService) AddNodesFromWithNodes(nodes []int64) {
	s.graphRepository.AddNodesFromWithNodes(nodes)
}

func (s *GraphService) AddNodesFromSourcesTargetsWeights(sources []int64, targets []int64, weights []float64) {
	s.graphRepository.AddNodesFromSourcesTargetsWeights(sources, targets, weights)
}

func (s *GraphService) GetTwitterAccountInfo(name string) int64 {
	return s.graphRepository.GetTwitterAccountInfo(name)
}
