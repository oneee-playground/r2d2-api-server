package dto

import "github.com/google/uuid"

type TaskListElem struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type TaskListOutput []TaskListElem

type TaskOutput struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Stage       string    `json:"stage"`
}

type TaskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskInput struct {
	IDInput
	TaskInput
}

type TaskStageInput struct {
	IDInput
	Stage string `json:"stage" validate:"task_stage"`
}
