package random

import "github.com/sethvargo/go-password/password"

func GenerateRandomPassword() (string, error) {
	password, err := password.Generate(15, 5, 5, false, false)
	if err != nil {
		return "", err
	}

	return password, nil
}
