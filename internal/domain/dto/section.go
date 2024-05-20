package dto

import "github.com/google/uuid"

type SectionListElem struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type SectionListOutput []SectionListElem

type SectionInput struct {
	Type        string `json:"type" validate:"section_type"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateSectionInput struct {
	IDInput
	SectionInput
}

type SectionIDInput struct {
	TaskID    uuid.UUID
	SectionID uuid.UUID
}

type UpdateSectionInput struct {
	SectionIDInput
	SectionInput
}

type SectionIndexInput struct {
	SectionIDInput
	Index int `json:"index"`
}
