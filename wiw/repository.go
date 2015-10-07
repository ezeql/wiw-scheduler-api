package wiw

//Repository is the services provider abstraction.
type Repository interface {
	ShiftsForUser(ID int) ([]Shift, error)
	ColleaguesForUser(ID int) ([]ColleaguesResult, error)
	ManagersForUser(ID int) ([]User, error)
	UserDetails(ID int) (User, error)
	CreateShift(shift *Shift) error
	UpdateOrCreateShift(shift *Shift) error
	ShiftsInRange(string, string) ([]Shift, error)
}
