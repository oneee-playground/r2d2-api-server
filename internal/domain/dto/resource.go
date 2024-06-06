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
	Image  string  `json:"image" binding:"required"`
	Name   string  `json:"name" binding:"required"`
	Port   uint16  `json:"port" binding:"required"`
	CPU    float64 `json:"cpu" binding:"required"`
	Memory uint64  `json:"memory" binding:"required"`
	// Boolean cannot be validated as a value. (e.g. value is false)
	// Reference: https://github.com/go-playground/validator/issues/142#issuecomment-127451987
	IsPrimary *bool `json:"isPrimary" binding:"required"`
}

type CreateResourceInput struct {
	IDInput
	ResourceInput
}

type ResourceIDInput struct {
	IDInput
	Name string `json:"string" binding:"required"`
}
