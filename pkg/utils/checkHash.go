package utils

import (
	"fmt"

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
