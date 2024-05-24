package resource_module

import (
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

func toResourceListOutput(resources []domain.Resource) *dto.ResourceListOutput {
	out := make(dto.ResourceListOutput, len(resources))

	for i, resource := range resources {
		out[i] = dto.ResourceListElem{
			Image:     resource.Image,
			Name:      resource.Name,
			Port:      resource.Port,
			CPU:       resource.CPU,
			Memory:    resource.Memory,
			IsPrimary: resource.IsPrimary,
		}
	}

	return &out
}
