package wiw

import (
	_ "github.com/go-sql-driver/mysql" //side effects for mysql
	"github.com/jinzhu/gorm"
	"strconv"
)

type MySQLRepository struct {
	gorm.DB
}

//NewMySQLRepo builds and initializes a MySQL db conncetion using passed dns arguments
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
		Select("ss.id AS shift_id, ss.employee_id AS colleague_id, ss.start_time AS colleague_start_time," +
		"ss.end_time AS colleague_end_time,s.start_time AS employee_start_time," +
		"s.end_time AS employee_end_time").
		Joins("INNER JOIN shifts AS ss ON s.employee_id = " + IDStr + " AND ss.employee_id !=" + IDStr).
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
	return user, result.Error
}

func (r *MySQLRepository) CreateShift(shift *Shift) error {
	return r.Create(shift).Error
}

func (r *MySQLRepository) UpdateOrCreateShift(shift *Shift) error {
	return r.Where(shift).FirstOrCreate(shift).Error
}
