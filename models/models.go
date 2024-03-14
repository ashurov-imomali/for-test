package models

import "time"

type CommissionRules struct {
	Id         int64      `gorm:"column:id" json:"id"`
	Index      int64      `gorm:"column:index" json:"index"`
	StartRange float64    `gorm:"column:start_range" json:"start_range"`
	EndRange   *float64   `gorm:"column:end_range" json:"end_range"`
	Value      float64    `gorm:"column:value" json:"value"`
	TypeId     int64      `gorm:"column:type_id" json:"type_id"`
	Active     *bool      `gorm:"active" json:"active"`
	ProfileId  int64      `gorm:"column:profile_id" json:"profile_id"`
	CreatedAt  *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt  *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CommissionRules) TableName() string {
	return "tcomission_rules"
}

type CommissionTypes struct {
	Id        int64      `gorm:"column:id" json:"id"`
	Name      string     `gorm:"column:name" json:"name"`
	Code      string     `gorm:"column:code" json:"code"`
	Active    bool       `gorm:"column:active" json:"active"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CommissionTypes) TableName() string {
	return "tcomission_types"
}

type CommissionProfiles struct {
	Id          int64      `gorm:"column:id" json:"id"`
	Name        string     `gorm:"column:name" json:"name"`
	Description string     `gorm:"column:description" json:"description"`
	CreatedBy   int64      `gorm:"column:created_by" json:"created_by"`
	UpdatedBy   int64      `gorm:"column:updated_by" json:"updated_by"`
	Active      bool       `gorm:"column:active" json:"active"`
	CreatedAt   *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (CommissionProfiles) TableName() string {
	return "tcomission_profiles"
}

type ProfileCreatRequest struct {
	Profile CommissionProfiles `json:"profile"`
	Rules   []CommissionRules  `json:"rules"`
}
