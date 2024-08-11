package utils

import (
	"fmt"
	"my_project/pkg/repository"
)

func verifyRole(role string) error {
	switch role {
	case repository.AdminRoleName:
		// Nothing to do, verified successfully.
	case repository.ModeratorRoleName:
		// Nothing to do, verified successfully.
	case repository.UserRoleName:
		// Nothing to do, verified successfully.
	default:
		// Return error message.
		return fmt.Errorf("role '%v' does not exist", role)
	}

	return nil
}
