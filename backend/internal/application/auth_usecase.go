package application

import (
	"fmt"
	"time"

	"app/internal/application/ports"
	"app/internal/domain/user"
)

type AuthUseCase struct {
	userRepo user.Repository
	verSvc   ports.VerificacionesService
}

func NewAuthUseCase(userRepo user.Repository, verSvc ports.VerificacionesService) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		verSvc:   verSvc,
	}
}

// Login выполняет логику /login
func (uc *AuthUseCase) Login(login, password string) (*user.User, error) {
	// 1. Try to find user by login
	usr, err := uc.userRepo.GetByLogin(login)
	if err != nil {
		// Assume if gorm.ErrRecordNotFound => create new
		// or handle error. Assume usr = nil
		// If error != not found, return err
		if !uc.userRepo.IsNotFoundError(err) {
			return nil, fmt.Errorf("repo error: %w", err)
		}
		usr = nil
	}

	// 2. If user NOT found OR user.ProviderID = 2 => check verificaciones
	if usr == nil || usr.ProviderID == 2 {
		// Call external verificaciones service
		resp, status, err := uc.verSvc.Login(ports.LoginReq{
			Username: login,
			Password: password,
		})
		if err != nil {
			return nil, fmt.Errorf("external login error: %w", err)
		}
		if status != 200 {
			return nil, fmt.Errorf("login failed with status %d", status)
		}

		if usr == nil {
			// User does not exist, create new
			newUser := &user.User{
				Login:       login,
				ProviderID:  2, // "verificaciones"
				CompanyID:   uint(resp.IdEmpresa),
				CompanyName: resp.Empresa,
				Active:      true,
				IsLogged:    true,
				// Password not stored (or empty string)
				Password: nil,
			}
			newUser.LastAccess = time.Now().Format("2006-01-02 15:04:05")

			if err := uc.userRepo.Create(newUser); err != nil {
				return nil, fmt.Errorf("create user error: %w", err)
			}

			// Get new user
			usr, err = uc.userRepo.GetByLogin(login)
			if err != nil {
				return nil, fmt.Errorf("get user error: %w", err)
			}

			return usr, nil
		} else {
			// user already exists, just update lastAccess
			usr.LastAccess = time.Now().Format("2006-01-02 15:04:05")
			usr.IsLogged = true
			if err := uc.userRepo.Update(usr); err != nil {
				return nil, fmt.Errorf("update user error: %w", err)
			}
			return usr, nil
		}
	}

	// 3. Otherwise, user found and ProviderID=1 => check local password
	if !usr.CheckPassword(password) {
		return nil, fmt.Errorf("invalid password")
	}

	// Password is ok, update lastAccess
	usr.LastAccess = time.Now().Format("2006-01-02 15:04:05")
	usr.IsLogged = true
	if err := uc.userRepo.Update(usr); err != nil {
		return nil, fmt.Errorf("update user error: %w", err)
	}
	return usr, nil
}
