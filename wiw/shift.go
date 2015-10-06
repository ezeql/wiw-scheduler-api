package wiw

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/now"
	"time"
)

//used for calculating hours worked in a whole week.
const HoursPerWeek = 24 * 7
const WorkHoursPerDay = 8

type Shift struct {
	Model
	Manager    *User     `json:"manager,omitempty"`
	ManagerID  uint      `json:"manager_id"`
	Employee   *User     `json:"employee,omitempty"`
	EmployeeID uint      `json:"employee_id,binding:"required""`
	BreakTime  float64   `json:"break",binding:"required"` //since go uses float64 for json
	StartTime  time.Time `json:"start_time",binding:"required",sql:"not null"`
	EndTime    time.Time `json:"end_time",binding:"required",sql:"not null"`
}

func (s *Shift) UnmarshalJSON(b []byte) error {
	const rfc2822Layout = "Mon Jan 02 15:04:05 -0700 2006"
	var parsed map[string]interface{}

	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	s.ManagerID = uint(parsed["manager_id"].(float64))
	s.EmployeeID = uint(parsed["employee_id"].(float64))

	s.BreakTime = parsed["break"].(float64)

	auxTime, err := time.Parse(rfc2822Layout, parsed["start_time"].(string))
	if err != nil {
		return err
	}
	s.StartTime = auxTime.UTC()

	auxTime, err = time.Parse(rfc2822Layout, parsed["end_time"].(string))
	if err != nil {
		return err
	}
	s.EndTime = auxTime.UTC()
	return nil
}

func (s Shift) String() string {
	return fmt.Sprintf("ID: %v  ManagerID: %v  EmployeeID: %v", s.ID, s.ManagerID, s.EmployeeID)
}

func SummarizeShifts(shifts []Shift) map[int][]float64 {
	now.FirstDayMonday = true
	summary := make(map[int][]float64)

	for _, value := range shifts {
		current := value.StartTime

		//advance until beggining of next week
		next := now.New(current).BeginningOfWeek().AddDate(0, 0, 7)
		isoYear, isoWeek := current.ISOWeek()

		//build key if not present
		buildKey(summary, isoYear, current)

		//set the first chunk of hours, until next week
		summary[isoYear][isoWeek-1] += next.Sub(current).Hours()

		current = next

		//iterate until end of shift
		for current.Before(value.EndTime) {
			isoYear, isoWeek = current.ISOWeek()

			buildKey(summary, isoYear, current)

			//note: break is part of work time.
			next = current.Add(time.Hour * HoursPerWeek)
			summary[isoYear][isoWeek-1] += next.Sub(current).Hours()
			current = next
		}
		//correct last iteration whole week time, better complexity than having the if inside the loop
		summary[isoYear][isoWeek-1] -= current.Sub(value.EndTime).Hours()
	}
	return summary
}

//builds key k with type []float64 with the required capacity for the year(52 or 53)
func buildKey(m map[int][]float64, k int, t time.Time) {
	v := ISOWeeksCount(t)
	if _, exists := m[k]; !exists {
		m[k] = make([]float64, v)
	}
}
