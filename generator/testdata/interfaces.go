package testdata

// User is a interface
type User interface {
	IsUser()
}

type Admin struct {
	ID        string
	Name      string
	Suspended bool
}

func (Admin) IsUser() {}

type Customer struct {
	ID    string
	Name  string
	Email string
}

func (Customer) IsUser() {}
