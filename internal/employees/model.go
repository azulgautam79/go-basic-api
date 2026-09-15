package employees

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	GenderMale           Gender = "male"
	GenderFemale         Gender = "female"
	GenderOther          Gender = "other"
	GenderPreferNotToSay Gender = "prefer_not_to_say"
)

type Employee struct {
	ID        uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	FirstName string     `json:"first_name" example:"John"`
	LastName  string     `json:"last_name" example:"Doe"`
	Email     string     `json:"email" example:"john.doe@example.com"`
	Gender    Gender     `json:"gender" example:"male"`
	IsActive  bool       `json:"is_active" example:"true"`
	CreatedAt time.Time  `json:"created_at" example:"2026-01-15T10:00:00Z"`
	UpdatedAt time.Time  `json:"updated_at" example:"2026-01-15T10:00:00Z"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" example:"null"`
}

type FilterParams struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Search   string `json:"search"`
	Gender   Gender `json:"gender"`
	IsActive *bool  `json:"is_active"`
}

type CreateEmployeeDTO struct {
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	Gender    Gender `json:"gender" example:"male" enums:"male,female,other,prefer_not_to_say"`
	IsActive  *bool  `json:"is_active,omitempty" example:"true"`
}

type UpdateEmployeeDTO struct {
	FirstName *string `json:"first_name,omitempty" example:"Jane"`
	LastName  *string `json:"last_name,omitempty" example:"Smith"`
	Email     *string `json:"email,omitempty" example:"jane.smith@example.com"`
	Gender    *Gender `json:"gender,omitempty" example:"female" enums:"male,female,other,prefer_not_to_say"`
	IsActive  *bool   `json:"is_active,omitempty" example:"false"`
}

type BatchDeleteDTO struct {
	IDs []uuid.UUID `json:"ids" example:"[\"550e8400-e29b-41d4-a716-446655440000\",\"a1b2c3d4-e5f6-7890-abcd-ef1234567890\"]"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page" example:"1"`
	Limit      int         `json:"limit" example:"10"`
	Total      int         `json:"total" example:"42"`
	TotalPages int         `json:"total_pages" example:"5"`
}
