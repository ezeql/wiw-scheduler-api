package wiw

import (
	"fmt"
)

type User struct {
	Model
	Name  string `json:"name" sql:"not null"`
	Role  string `json:"role" sql:"not null"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func (b User) String() string {
	return fmt.Sprintf("ID: %v  Name: %v  Role: %v  Email: %v  Phone: %v", b.ID, b.Name, b.Role, b.Email, b.Phone)
}
