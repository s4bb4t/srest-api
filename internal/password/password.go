package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	const op = "password.HashPassword"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	fmt.Printf(string(hashedPassword))

	return hashedPassword, nil
}

// func CheckPassword(hashedPassword []byte, password string) error {
// 	const op = "password.CheckPassword"

// 	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
// 	if err != nil {
// 		return fmt.Errorf("%s: %v", op, err)
// 	}

// 	return nil
// }
