package verificaciones

import (
	"app/internal/application/ports"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)

func (vc *verificacionesClient) Login(req ports.LoginReq) (*ports.LoginRes, int, error) {
	resp, err := vc.client.R().
		SetQueryParams(map[string]string{
			"username": req.Username,
			"password": req.Password,
		}).
		Post(vc.baseURL + vc.loginRoute)
	if err != nil {
		// Error when requesting (no response, network failure, etc.)
		return nil, http.StatusInternalServerError, NewServiceError(ErrLoginRequestFailed, err)
	}

	if resp.StatusCode() != http.StatusOK {
		// Service returned not 200 => treat as ErrLoginFailed
		return nil, resp.StatusCode(), NewServiceError(ErrLoginFailed, errors.New(resp.Status()))
	}

	// Example security check
	var securityResp struct {
		Security string `json:"security"`
	}
	if err := json.Unmarshal(resp.Body(), &securityResp); err == nil && securityResp.Security == "failed" {
		return nil, http.StatusUnauthorized, NewServiceError(ErrSecurityFailed, nil)
	}

	// Parse response as ports.LoginRes
	var response ports.LoginRes
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, resp.StatusCode(), NewServiceError(ErrParseLoginResponse, err)
	}

	if response.Token == "" || response.AppToken == "" {
		return nil, resp.StatusCode(), NewServiceError(ErrInvalidLoginResponse, nil)
	}

	// If you need to pull the company name immediately
	company, err := vc.GetCompanyByCompanyId(strconv.Itoa(response.IdEmpresa))
	if err != nil {
		return nil, resp.StatusCode(), NewServiceError(ErrGetCompanyFailed, err)
	}
	if company != nil {
		response.Empresa = company.CompanyName
	}

	return &response, http.StatusOK, nil
}

// --- systemLogin ---
// Used by methods GetCompanyByCompanyId, GetCompanyByICCID, CheckIfUserExists,
// if they need an appToken
func (vc *verificacionesClient) systemLogin() (string, error) {
	// If the token is still in the cache
	if vc.token != "" && time.Now().Before(vc.tokenExp) {
		return vc.token, nil
	}

	// Need to login to the system
	resp, err := vc.client.R().
		SetQueryParams(map[string]string{
			"username": vc.systemUsername,
			"password": vc.systemPassword,
		}).
		Post(vc.baseURL + vc.loginRoute)
	if err != nil {
		return "", NewServiceError(ErrLoginRequestFailed, err)
	}
	if resp.StatusCode() != http.StatusOK {
		return "", NewServiceError(ErrLoginFailed, errors.New(resp.Status()))
	}

	var lr ports.LoginRes
	if err := json.Unmarshal(resp.Body(), &lr); err != nil {
		return "", NewServiceError(ErrParseLoginResponse, err)
	}

	if lr.AppToken == "" {
		return "", NewServiceError(ErrEmptyTokenResponse, nil)
	}

	vc.token = lr.AppToken
	vc.tokenExp = time.Now().Add(vc.tokenCacheDuration)
	return vc.token, nil
}

// --- GetCompanyByCompanyId ---

func (vc *verificacionesClient) GetCompanyByCompanyId(companyId string) (*ports.GetCompanyByCompanyIdRes, error) {
	token, err := vc.systemLogin()
	if err != nil {
		return nil, NewServiceError(ErrTokenFailed, err)
	}

	resp, err := vc.client.R().
		SetQueryParams(map[string]string{
			"apptoken":  token,
			"companyId": companyId,
		}).
		Post(vc.baseURL + vc.getCompanyByCompanyIdRoute)
	if err != nil {
		return nil, NewServiceError(ErrCompanyRequestFailed, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, NewServiceError(ErrCompanyRequestFailed, errors.New(resp.Status()))
	}

	var response ports.GetCompanyByCompanyIdRes
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, NewServiceError(ErrParseCompanyResponse, err)
	}

	return &response, nil
}

// --- GetCompanyByICCID ---

func (vc *verificacionesClient) GetCompanyByICCID(iccid string) (*ports.GetCompanyByICCIDRes, error) {
	token, err := vc.systemLogin()
	if err != nil {
		return nil, NewServiceError(ErrTokenFailed, err)
	}

	resp, err := vc.client.R().
		SetQueryParams(map[string]string{
			"apptoken": token,
			"iccid":    iccid,
		}).
		Post(vc.baseURL + vc.getCompanyByICCIDRoute)
	if err != nil {
		return nil, NewServiceError(ErrICCIDRequestFailed, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, NewServiceError(ErrICCIDRequestFailed, errors.New(resp.Status()))
	}

	var response ports.GetCompanyByICCIDRes
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, NewServiceError(ErrParseICCIDResponse, err)
	}

	return &response, nil
}

// --- CheckIfUserExists ---

func (vc *verificacionesClient) CheckIfUserExists(username string) (bool, error) {
	token, err := vc.systemLogin()
	if err != nil {
		return false, NewServiceError(ErrTokenFailed, err)
	}

	resp, err := vc.client.R().
		SetQueryParams(map[string]string{
			"apptoken": token,
			"tech":     username,
		}).
		Post(vc.baseURL + vc.checkUserExistsRoute)
	if err != nil {
		return false, NewServiceError(ErrCheckUserExistsFailed, err)
	}

	if resp.StatusCode() != http.StatusOK {
		return false, NewServiceError(ErrCheckUserExistsFailed, errors.New(resp.Status()))
	}

	var response ports.CheckIfUserExistsRes
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return false, NewServiceError(ErrParseUserExistsResponse, err)
	}

	userExists := false

	if strRol, ok := response.Rol.(string); ok {
		userExists = (strRol == "3" && response.HasLogin == "true")
	} else if floatRol, ok := response.Rol.(float64); ok {
		userExists = (int(floatRol) == 3 && response.HasLogin == "true")
	}

	return userExists, nil
}
