package application

import (
	"fmt"
	"time"

	"app/internal/application/ports"
	"app/internal/domain/role"
	"app/internal/domain/user"
	"app/internal/infrastructure/token/paseto"
	"app/internal/infrastructure/token/refresh"
)

type AuthUseCase struct {
	userRepo    user.Repository
	verSvc      ports.VerificacionesService
	roleRepo    role.RoleRepository
	userService *user.UserService
}

func NewAuthUseCase(userRepo user.Repository, verSvc ports.VerificacionesService, roleRepo role.RoleRepository) *AuthUseCase {
	userService := user.NewUserService(userRepo, roleRepo)
	return &AuthUseCase{
		userRepo:    userRepo,
		verSvc:      verSvc,
		roleRepo:    roleRepo,
		userService: userService,
	}
}

func (uc *AuthUseCase) Login(login, password string) (*user.User, error) {
	// 1. Try to find user by login
	usr, err := uc.userRepo.GetByLogin(login)
	if err != nil {
		if !uc.userRepo.IsNotFoundError(err) {
			// another error
			return nil, fmt.Errorf("repo error: %w", err)
		}
		usr = nil
	}

	// 2. If user NOT found OR user.ProviderID=2 => check external service
	if usr == nil || usr.ProviderID == 2 {
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
			// Create new user
			newUser := &user.User{
				Login:        login,
				ProviderID:   2,
				ProviderName: "Verificaciones",
				CompanyID:    uint(resp.IdEmpresa),
				CompanyName:  resp.Empresa,
				Active:       true,
				IsLogged:     true,
				Password:     nil, // do not store password
			}
			newUser.LastAccess = time.Now().Format("2006-01-02 15:04:05")

			if err := uc.userRepo.Create(newUser); err != nil {
				return nil, fmt.Errorf("create user error: %w", err)
			}
			// After creation, get user from DB again
			usr, err = uc.userRepo.GetByLogin(login)
			if err != nil {
				return nil, fmt.Errorf("get user error: %w", err)
			}
		} else {
			// user already exists
			usr.LastAccess = time.Now().Format("2006-01-02 15:04:05")
			usr.IsLogged = true
			if err := uc.userRepo.Update(usr); err != nil {
				return nil, fmt.Errorf("update user error: %w", err)
			}
		}
	} else {
		// 3. If user found + ProviderID=1 => local password check
		if !usr.CheckPassword(password) {
			return nil, fmt.Errorf("invalid password")
		}
		usr.LastAccess = time.Now().Format("2006-01-02 15:04:05")
		usr.IsLogged = true
		if err := uc.userRepo.Update(usr); err != nil {
			return nil, fmt.Errorf("update user error: %w", err)
		}
	}

	// 4. Check if `OwnerID` exists, if yes, get owner
	var ownerUsername string
	if usr.OwnerID != nil {
		ownerUser, err := uc.userRepo.GetByID(*usr.OwnerID)
		if err == nil {
			ownerUsername = ownerUser.Login
		}
	}

	// At this point, usr is definitely not nil and is authorized
	// 5. Generate refresh-token + expDate
	token, expDate, err := refresh.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("refresh token generation error: %w", err)
	}
	usr.Refresh = &token
	usr.RefreshExp = expDate

	// 6. Save refresh-token in DB
	if err := uc.userRepo.Update(usr); err != nil {
		return nil, fmt.Errorf("update user (refresh) error: %w", err)
	}

	// 7. Check and assign roles to user
	if err := uc.userService.EnsureUserRoles(usr); err != nil {
		return nil, fmt.Errorf("ensure user roles error: %w", err)
	}

	// 8. Get user roles directly from database
	userRoles, err := uc.userService.GetUserRoles(usr.ID)
	if err != nil {
		return nil, fmt.Errorf("get user roles error: %w", err)
	}

	// Assign roles to user object
	usr.Roles = userRoles

	// 9. Generate access-token
	// Convert roles to string for token
	roleNames := uc.userService.GetRoleNamesString(userRoles)

	accessToken, _, err := paseto.Paseto().GenerateToken(paseto.PasetoClaims{
		Username:      usr.Login,
		CompanyID:     int(usr.CompanyID),
		CompanyName:   usr.CompanyName,
		Roles:         roleNames,
		OwnerUsername: ownerUsername,
		// IsPrimary:     usr.Profile != nil && usr.Profile.IsPrimary,
	})
	if err != nil {
		return nil, fmt.Errorf("access token generation error: %w", err)
	}
	usr.AccessToken = accessToken
	return usr, nil
}
