package entity

import "time"

type TemplateTask struct {
	TetaID       int64     `gorm:"primaryKey;column:teta_id"`
	TetaName     string    `gorm:"column:teta_name"`
	TetaRoleID   int64     `gorm:"column:teta_role_id"`
	Description  string    `gorm:"column:description"`
	SequenceNo   int16     `gorm:"column:sequence_no"`
	IsRequired   bool      `gorm:"column:is_required"`
	CallFunction *string   `gorm:"column:call_function"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`

	// Relationships
	Attributes []TemplateTaskAttribute `gorm:"foreignKey:TetatTetaID;references:TetaID"`
}

func (TemplateTask) TableName() string {
	return "mst.template_tasks"
}

type TemplateTaskAttribute struct {
	TetatID       int64     `gorm:"primaryKey;column:tetat_id"`
	TetatName     string    `gorm:"column:tetat_name"`
	TetatTetaID   int64     `gorm:"column:tetat_teta_id"`
	Description   string    `gorm:"column:description"`
	IsRequired    bool      `gorm:"column:is_required"`
	AttributeType string    `gorm:"column:attribute_type"`
	CreatedAt     time.Time `gorm:"column:created_at"`
}

func (TemplateTaskAttribute) TableName() string {
	return "mst.template_task_attributes"
}

type LeasingTask struct {
	TaskID          int64      `gorm:"primaryKey;column:task_id"`
	TaskName        string     `gorm:"column:task_name"`
	TemplateTaskID  *int64     `gorm:"column:template_task_id"`
	Startdate       *time.Time `gorm:"column:startdate"`
	Enddate         *time.Time `gorm:"column:enddate"`
	ActualStartdate *time.Time `gorm:"column:actual_startdate"`
	ActualEnddate   *time.Time `gorm:"column:actual_enddate"`
	Status          string     `gorm:"column:status"`
	ContractID      int64      `gorm:"column:contract_id"`
	RoleID          int64      `gorm:"column:role_id"`
	SequenceNo      int16      `gorm:"column:sequence_no"`
	CallFunction    *string    `gorm:"column:call_function"`
	CreatedAt       time.Time  `gorm:"column:created_at"`

	// Relationships
	Attributes []LeasingTaskAttribute `gorm:"foreignKey:TasaTaskID;references:TaskID"`
}

func (LeasingTask) TableName() string {
	return "leasing.leasing_tasks"
}

type LeasingTaskAttribute struct {
	TasaID     int64     `gorm:"primaryKey;column:tasa_id"`
	TasaName   string    `gorm:"column:tasa_name"`
	TasaValue  *string   `gorm:"column:tasa_value"`
	TasaStatus string    `gorm:"column:tasa_status"`
	TasaTaskID int64     `gorm:"column:tasa_task_id"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (LeasingTaskAttribute) TableName() string {
	return "leasing.leasing_tasks_attributes"
}
