package testdata

import (
	"fmt"
	"io"
	"strconv"
)

type User struct {
	ID   *string
	Role *Role
}

// Role a user role
type Role string

const (
	Role_ROOT       Role = "ROOT"
	Role_SUPERVISOR Role = "SUPERVISOR"
	Role_USER       Role = "USER"
	Role_ANONYMOUS  Role = "ANONYMOUS"
)

var AllRole = []Role{
	Role_ROOT,
	Role_SUPERVISOR,
	Role_USER,
	Role_ANONYMOUS,
}

func (e Role) IsValid() bool {
	switch e {
	case Role_ROOT, Role_SUPERVISOR, Role_USER, Role_ANONYMOUS:
		return true
	}
	return false
}

func (e Role) String() string {
	return string(e)
}

func (e *Role) UnmarshalJSON(v []byte) error {
	*e = Role(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Role", str)
	}
	return nil
}

func (e Role) MarshalJSON(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
