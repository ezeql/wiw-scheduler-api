package wiw

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/now"
	"strconv"
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
	EmployeeID uint      `json:"employee_id" binding:"required"`
	BreakTime  float64   `json:"break" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required" sql:"not null"`
	EndTime    time.Time `json:"end_time" binding:"required" sql:"not null"`
}

//workaround since gorm wont recognize time.Time composed types
func (s *Shift) UnmarshalJSON(b []byte) error {
	const rfc2822Layout = "Mon Jan 02 15:04:05 -0700 2006"
	var parsed map[string]interface{}

	if err := json.Unmarshal(b, &parsed); err != nil {
		return err
	}

	if parsed["start_time"] == nil {
		return errors.New("missing start time")
	}

	auxTime, err := time.Parse(rfc2822Layout, parsed["start_time"].(string))
	if err != nil {
		return errors.New("invalid start time")
	}
	s.StartTime = auxTime.UTC()

	if parsed["end_time"] == nil {
		return errors.New("missing end time")
	}

	auxTime, err = time.Parse(rfc2822Layout, parsed["end_time"].(string))
	if err != nil {
		return errors.New("invalid end time")
	}
	s.EndTime = auxTime.UTC()

	if parsed["manager_id"] != nil {
		id, ok := parsed["manager_id"].(float64)
		if !ok {
			return errors.New("invalid manager_id")
		}
		s.ManagerID = uint(id)
	}

	if parsed["employee_id"] != nil {
		id, ok := parsed["employee_id"].(float64)
		if !ok {
			return errors.New("invalid employee_id")
		}
		s.EmployeeID = uint(id)
	}

	if parsed["break"] == nil {

	}

	return nil
}

func (s Shift) String() string {
	return fmt.Sprintf("ID: %v  ManagerID: %v  EmployeeID: %v", s.ID, s.ManagerID, s.EmployeeID)
}

func SummarizeShifts(shifts []Shift) map[string][]float64 {
	now.FirstDayMonday = true
	summary := make(map[string][]float64)

	for _, value := range shifts {
		current := value.StartTime

		//advance until beggining of next week
		next := now.New(current).BeginningOfWeek().AddDate(0, 0, 7)
		isoYear, isoWeek := current.ISOWeek()
		isoYearStr := strconv.Itoa(isoYear)

		//build key if not present
		buildKey(summary, isoYear)

		//set the first chunk of hours, until next week
		summary[isoYearStr][isoWeek-1] += next.Sub(current).Hours()

		current = next

		//iterate until end of shift
		for current.Before(value.EndTime) {
			isoYear, isoWeek = current.ISOWeek()
			isoYearStr = strconv.Itoa(isoYear)
			buildKey(summary, isoYear)

			//note: break is part of work time.
			next = current.Add(time.Hour * HoursPerWeek)
			summary[isoYearStr][isoWeek-1] += next.Sub(current).Hours()
			current = next
		}
		//correct last iteration whole week time, better complexity than having the if inside the loop
		summary[isoYearStr][isoWeek-1] -= current.Sub(value.EndTime).Hours()
	}
	return summary
}

//builds key k with type []float64 with the required capacity for the year(52 or 53)
func buildKey(m map[string][]float64, y int) {
	v := ISOWeeksCount(y)
	k := strconv.Itoa(y)
	if _, exists := m[k]; !exists {
		m[k] = make([]float64, v)
	}
}
