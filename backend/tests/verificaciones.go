package tests

import (
	"app/internal/application/ports"
	"app/internal/infrastructure/webhooks/verificaciones"
	"app/pkg/config"
	"fmt"
	"net/http"
)

func TestLoginSuccess() {

	// 1. Get username/password
	username := config.ENV().VERIFICACIONES_USERNAME
	password := config.ENV().VERIFICACIONES_PASSWORD

	// 2. Create verSvc (verificaciones service)
	verSvc := verificaciones.NewVerificacionesClient()

	// 3. Prepare the request (struct from ports)
	req := ports.LoginReq{
		Username: username,
		Password: password,
	}

	// 4. Call the method (ports)
	loginRes, statusCode, err := verSvc.Login(req)
	if err != nil {
		// Error (maybe *ServiceError)
		fmt.Printf("Login failed: %v\n", err)
		return
	}

	if statusCode != http.StatusOK {
		fmt.Printf("Login returned unexpected code: %d\n", statusCode)
		return
	}

	// 5. If everything is ok, print the result
	fmt.Println("Login successful. Tokens:")
	fmt.Printf(" token=%s, appToken=%s\n", loginRes.Token, loginRes.AppToken)
}

func TestLoginFailed() {

	// 1. Get username/password
	username := "pepe"
	password := "pepino"

	// 2. Create verSvc (verificaciones service)
	verSvc := verificaciones.NewVerificacionesClient()

	// 3. Prepare the request (struct from ports)
	req := ports.LoginReq{
		Username: username,
		Password: password,
	}

	// 4. Call the method (ports)
	loginRes, statusCode, err := verSvc.Login(req)
	if err != nil {
		// Error (maybe *ServiceError)
		fmt.Printf("Login failed: %v\n", err)
		return
	}

	if statusCode != http.StatusOK {
		fmt.Printf("Login returned unexpected code: %d\n", statusCode)
		return
	}

	// 5. If everything is ok, print the result
	fmt.Println("Login successful. Tokens:")
	fmt.Printf(" token=%s, appToken=%s\n", loginRes.Token, loginRes.AppToken)
}
