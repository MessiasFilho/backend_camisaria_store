package schemas

import "time"

type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleUser   UserRole = "user"
	RoleClient UserRole = "client"
)

type Users struct {
	ID        uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"password" gorm:"type:varchar(255)"`
	Role      UserRole  `json:"role" gorm:"type:enum('admin','user','client');default:'client'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
