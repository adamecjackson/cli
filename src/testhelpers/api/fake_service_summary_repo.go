package api

import "cf/models"

type FakeServiceSummaryRepo struct {
	GetSummariesInCurrentSpaceInstances []models.ServiceInstance
}

func (repo *FakeServiceSummaryRepo) GetSummariesInCurrentSpace() (instances []models.ServiceInstance, apiErr error) {
	instances = repo.GetSummariesInCurrentSpaceInstances
	return
}
