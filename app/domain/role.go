package domain

type Role string

const (
	UserRole       Role = "user"
	StaffRole      Role = "staff"
	BackofficeRole Role = "backoffice"
	SystemRole     Role = "system"
)

func (r Role) IsBackoffice() bool {
	return r == BackofficeRole
}

func (r Role) IsUser() bool {
	return r == UserRole
}

func (r Role) IsStaff() bool {
	return r == StaffRole
}

func (r Role) IsSystem() bool {
	return r == SystemRole
}

func (r Role) String() string {
	return string(r)
}
