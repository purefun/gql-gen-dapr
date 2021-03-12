package testdata

// User a union
type User interface {
	IsUser()
}

type Admin struct {
	ID string
}

func (Admin) IsUser() {}

type Customer struct {
	ID string
}

func (Customer) IsUser() {}
