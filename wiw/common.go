package wiw

import "time"

type Model struct {
	ID        uint      `gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type ColleaguesResult struct {
	ShiftID            int
	ColleagueID        int
	ColleagueStartTime time.Time
	ColleagueEndTime   time.Time
	EmployeeStartTime  time.Time
	EmployeeEndTime    time.Time
}

type UserShift struct {
	User
	Shift
}
