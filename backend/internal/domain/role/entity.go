package role

type Role struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	Desc string `json:"desc"`
}

type RefRoleUser struct {
	UserID uint `json:"userId"`
	RoleID uint `json:"roleId"`
}
