package dto

import "github.com/google/uuid"

type SectionListElem struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	RPM         uint64    `json:"rpm"`
	Example     string    `json:"example"`
}

type SectionListOutput []SectionListElem

type SectionInput struct {
	Type        string `json:"type" binding:"required" validate:"section_type"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	RPM         uint64 `json:"rpm" binding:"required"`
}

type CreateSectionInput struct {
	IDInput
	SectionInput
}

type SectionIDInput struct {
	TaskID    uuid.UUID `uri:"taskID" binding:"required"`
	SectionID uuid.UUID `uri:"sectionID" binding:"required"`
}

type UpdateSectionInput struct {
	SectionIDInput
	SectionInput
}

type SectionIndexInput struct {
	SectionIDInput
	Index int `json:"index" binding:"required"`
}
