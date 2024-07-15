package dto

type SectionListElem struct {
	ID          string `json:"id" binding:"uuid"`
	Type        string `json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	RPM         uint64 `json:"rpm"`
	Example     string `json:"example"`
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
	TaskID    string `uri:"id" binding:"required,uuid"`
	SectionID string `uri:"sectionID" binding:"required,uuid"`
}

type UpdateSectionInput struct {
	SectionIDInput
	SectionInput
}

type SectionIndexInput struct {
	SectionIDInput
	Index int `json:"index" binding:"required"`
}
