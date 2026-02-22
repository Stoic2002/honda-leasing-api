package entity

import "time"

type Role struct {
	RoleID      int64     `gorm:"primaryKey;column:role_id"`
	RoleName    string    `gorm:"column:role_name"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (Role) TableName() string {
	return "account.roles"
}

type User struct {
	UserID         int64      `gorm:"primaryKey;column:user_id"`
	Username       string     `gorm:"column:username"`
	PhoneNumber    string     `gorm:"column:phone_number"`
	Email          string     `gorm:"column:email"`
	FullName       string     `gorm:"column:full_name"`
	Password       string     `gorm:"column:password"`
	PinKey         *string    `gorm:"column:pin_key"`
	IsActive       bool       `gorm:"column:is_active"`
	LastLogin      *time.Time `gorm:"column:last_login"`
	FailedAttempts int16      `gorm:"column:failed_attempts"`
	LockedUntil    *time.Time `gorm:"column:locked_until"`
	CreatedAt      time.Time  `gorm:"column:created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at"`
	Roles          []Role     `gorm:"many2many:account.user_roles;joinForeignKey:user_id;joinReferences:role_id"`
}

func (User) TableName() string {
	return "account.users"
}
