package utils

import (
	"fmt"
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/modules/users"
	"golang.org/x/crypto/bcrypt"
)

func BcryptHashing(obj *users.UserRegisterReq) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hash password failed: %v", err)
	}
	obj.Password = string(hashPassword)
	return nil
}

func CheckTaskStatus(status string) (string, error) {
	status = strings.ToUpper(status)

	switch status {
	case "TODO", "DOING", "DONE":
		return status, nil
	default:
		return "", fmt.Errorf("invalid task status")
	}
}
