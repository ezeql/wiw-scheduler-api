package wiw

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/now"
	"math"
	"time"
)

//used for calculating hours worked in a whole week.
const HoursPerWeek = 24 * 7

const WeeksPerYear = 53

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
	//s.EmployeeID = uint(parsed["employee_id"].(float64))

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

func SummarizeShifts(shifts []Shift) map[string][]float64 {
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
