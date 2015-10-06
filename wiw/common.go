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

func ISOWeeksCount(t time.Time) int {

	tt := time.Date(t.Year(), time.December, 31, 0, 0, 0, 0, t.Location())
	ordinalDay := tt.YearDay()
	weekDay := int(tt.Weekday()) - 1
	return (ordinalDay - weekDay + 10) / 7

}
