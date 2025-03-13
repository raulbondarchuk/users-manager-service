package application

import (
	"app/internal/application/ports"
	"app/internal/domain/role"
	"app/internal/domain/user"
	"fmt"
	"net/http"
	"time"
)

const PROVIDER_SECONDARY = 3

type SubUserUseCase struct {
	userRepo          user.Repository
	roleRepo          role.RoleRepository
	userService       *user.UserService
	verificacionesSvc ports.VerificacionesService
}

func NewSubUserUseCase(userRepo user.Repository, roleRepo role.RoleRepository, verificacionesSvc ports.VerificacionesService) *SubUserUseCase {
	return &SubUserUseCase{
		userRepo:          userRepo,
		roleRepo:          roleRepo,
		userService:       user.NewUserService(userRepo, roleRepo),
		verificacionesSvc: verificacionesSvc,
	}
}

// CreateSubUser creates a subuser for a given main user's username
func (uc *SubUserUseCase) CreateSubUser(mainUsername, subUsername, subPassword, roles, email string) (*user.User, error) {

	// Check if user exists in verificaciones
	exists, err := uc.verificacionesSvc.CheckIfUserExists(subUsername)
	if err != nil {
		return nil, fmt.Errorf("error checking if user exists in verificaciones: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("user %s already exists in verificaciones", subUsername)
	}

	// Start a new transaction
	tx := uc.userRepo.BeginTransaction()

	// Ensure transaction is rolled back in case of error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Println("Transaction rolled back due to panic:", r)
		}
	}()

	// 1. Get main user by username
	mainUser, err := uc.userRepo.GetByLogin(mainUsername)
	if err != nil {
		tx.Rollback()
		if uc.userRepo.IsNotFoundError(err) {
			return nil, fmt.Errorf("main user not found: %s", mainUsername)
		}
		return nil, fmt.Errorf("error retrieving main user: %w", err)
	}

	// 2. Create subuser
	subUser := &user.User{
		Login:       subUsername,
		OwnerID:     &mainUser.ID,
		Active:      true,
		IsLogged:    false,
		CompanyID:   mainUser.CompanyID,
		CompanyName: mainUser.CompanyName,

		ProviderID: uint(PROVIDER_SECONDARY),
		CreatedAt:  time.Now().Format(time.RFC3339),
		LastAccess: time.Now().Format(time.RFC3339),
	}

	// Set password for subuser
	if err := subUser.SetPassword(subPassword); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error setting password for subuser: %w", err)
	}

	// 3. Save subuser to the repository
	if err := uc.userRepo.CreateWithTransaction(tx, subUser); err != nil {
		tx.Rollback()
		if uc.userRepo.IsAlreadyExistsError(err) {
			return nil, fmt.Errorf("user already exists")
		}
		return nil, fmt.Errorf("error creating subuser: %w", err)
	}

	subUser, err = uc.userRepo.GetByLoginWithTransaction(tx, subUsername)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error retrieving subuser: %w", err)
	}

	if roles != "" {
		// 4. Assign additional roles to subuser
		fmt.Println("roles", roles)
		if err := uc.userService.AssignRolesToSubUser(tx, subUser, roles); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error assigning roles to subuser: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	// 5. Ensure subuser has the necessary roles
	if err := uc.userService.EnsureUserRoles(subUser); err != nil {
		return nil, fmt.Errorf("error ensuring roles for subuser: %w", err)
	}

	subUser, err = uc.userRepo.GetByLogin(subUsername)
	if err != nil {
		return nil, fmt.Errorf("error retrieving subuser: %w", err)
	}

	return subUser, nil
}

func (uc *SubUserUseCase) DeleteSubuser(username string, companyId uint) (int, error) {

	user, err := uc.userRepo.GetByLogin(username)
	if err != nil {
		if uc.userRepo.IsNotFoundError(err) {
			return http.StatusNotFound, err
		}
		return http.StatusInternalServerError, err
	}

	if user.CompanyID != companyId || user.OwnerID == nil {
		return http.StatusNotFound, fmt.Errorf("user is not a subuser of this company")
	}

	return http.StatusOK, uc.userRepo.DeleteUserByUsername(username)
}
