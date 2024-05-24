package dto

type ResourceListElem struct {
	Image     string  `json:"image"`
	Name      string  `json:"name"`
	Port      uint16  `json:"port"`
	CPU       float64 `json:"cpu"`
	Memory    uint64  `json:"memory"`
	IsPrimary bool    `json:"isPrimary"`
}

type ResourceListOutput []ResourceListElem

type ResourceInput struct {
	Image     string  `json:"image"`
	Name      string  `json:"name"`
	Port      uint16  `json:"port"`
	CPU       float64 `json:"cpu"`
	Memory    uint64  `json:"memory"`
	IsPrimary bool    `json:"isPrimary"`
}

type CreateResourceInput struct {
	IDInput
	ResourceInput
}

type ResourceIDInput struct {
	IDInput
	Name string `json:"string"`
}
