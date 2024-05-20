package section_module

import (
	"github.com/oneee-playground/r2d2-api-server/internal/domain"
	"github.com/oneee-playground/r2d2-api-server/internal/domain/dto"
)

func toSectionListOutput(sections []domain.Section) *dto.SectionListOutput {
	out := make(dto.SectionListOutput, len(sections))
	for i, section := range sections {
		out[i] = dto.SectionListElem{
			ID:          section.ID,
			Type:        string(section.Type),
			Title:       section.Title,
			Description: section.Description,
		}
	}

	return &out
}
