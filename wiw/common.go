package wiw

import "time"

//Model default  used in all table mappings
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

func ISOWeeksCount(y int) int {
	tt := time.Date(y, time.December, 31, 0, 0, 0, 0, time.UTC)
	_, week := tt.ISOWeek()
	if week == 53 {
		return 53
	}
	return 52
}
