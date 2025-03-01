package ports

type VerificacionesService interface {
	Login(request LoginReq) (*LoginRes, int, error)
	GetCompanyByCompanyId(companyId string) (*GetCompanyByCompanyIdRes, error)
	GetCompanyByICCID(iccid string) (*GetCompanyByICCIDRes, error)
	CheckIfUserExists(username string) (bool, error)
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRes struct {
	Usuario     string `json:"usuario"`
	IdEmpresa   int    `json:"idempresa"`
	Empresa     string `json:"empresa"`
	Rol         int    `json:"rol"`
	TipoCliente string `json:"tipocliente"`
	Idioma      string `json:"idioma"`
	Token       string `json:"token"`
	Ce          int    `json:"ce"`
	AppToken    string `json:"apptoken"`
}

type Company struct {
	Codigocliente string `json:"codigocliente"`
	Nombrecliente string `json:"nombrecliente"`
}

type GetCompanyByCompanyIdRes struct {
	CompanyId   string `json:"companyId"`
	CompanyName string `json:"companyName"`
}

// Change the response type to a slice
type GetCompanyByICCIDRes []Company

type CheckIfUserExistsRes struct {
	TechExists  string      `json:"techexists"`
	HasLogin    string      `json:"haslogin"`
	Rol         interface{} `json:"rol"`
	IdCompany   interface{} `json:"idcompany"`
	NameCompany string      `json:"namecompany"`
}
