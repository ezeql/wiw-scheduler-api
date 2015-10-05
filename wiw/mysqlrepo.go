package wiw

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" //side effects
	"github.com/jinzhu/gorm"
	"strconv"
)

type MySQLRepository struct {
	gorm.DB
}

func NewMySQLRepo(dsn string) (*MySQLRepository, error) {
	var err error
	r := &MySQLRepository{}

	if r.DB, err = gorm.Open("mysql", dsn); err != nil {
		return nil, err
	}

	r.AutoMigrate(&User{}, &Shift{})
	r.LogMode(true)
	return r, r.Error
}

func (r *MySQLRepository) ShiftsForUser(ID int) ([]Shift, error) {

	user := User{}
	shifts := []Shift{}

	if result := r.First(&user, ID); result.Error != nil {
		return nil, result.Error
	}

	result := r.DB.Model(&user).Related(&shifts, "EmployeeID")
	return shifts, result.Error
}

func (r *MySQLRepository) ColleaguesForUser(ID int) ([]ColleaguesResult, error) {
	results := []ColleaguesResult{}
	IDStr := strconv.Itoa(ID)
	result := r.Table("shifts as s").
		Select("ss.id as shift_id, ss.employee_id as colleague_id, ss.start_time as colleague_start_time," +
		"ss.end_time as colleague_end_time,s.start_time as employee_start_time," +
		"s.end_time as employee_end_time").
		Joins("inner join shifts as ss on s.employee_id = " + IDStr + " AND ss.employee_id !=" + IDStr).
		Where("s.start_time < ss.end_time AND s.end_time > ss.start_time").
		Find(&results)

	return results, result.Error
}

func (r *MySQLRepository) ManagersForUser(ID int) ([]User, error) {
	managers := []User{}
	result := r.Table("users u, shifts s").Select("u.*").
		Where("u.id = s.manager_id  AND s.employee_id = ?", ID).Find(&managers)
	return managers, result.Error

}

func (r *MySQLRepository) ShiftsInRange(from string, to string) ([]Shift, error) {
	shifts := []Shift{}
	result := r.Where("start_time BETWEEN ? AND ?", from, to).Preload("Manager").Find(&shifts)
	return shifts, result.Error
}
func (r *MySQLRepository) UserDetails(ID int) (User, error) {
	user := User{}
	result := r.First(&user, ID)
	fmt.Println("%v", result.Error)

	return user, result.Error
}

func (r *MySQLRepository) UpdateShift() (Shift, error) {
	newShift := Shift{}
	// shift := &Shift{}
	//shiftID := c.Param("id")
	/*
		if err := c.BindJSON(&newShift); err == nil {
			if r.First(&shift, shiftID).RecordNotFound() {
				forcedID, _ := strconv.Atoi(shiftID)
				newShift.ID = uint(forcedID)
				r.Create(&newShift)
				c.JSON(200, &newShift)
			} else {
				r.Model(&shift).Updates(&newShift)
				c.JSON(200, &shift)
			}
		} else {
			c.JSON(http.StatusBadRequest, err.Error())
		}*/
	return newShift, nil
}
