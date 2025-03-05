package validation

// TODO: Realizar una func (expiredToken, refresh) para generar un nuevo par de tokens.
// Hay que encontrar usuario por refresh token, comprobar si ha expirado, generar un nuevo par de tokens y actualizar el usuario con los nuevos tokens.

// // CheckPairTokens checks if the jwtToken and refreshToken are valid
// func CheckPairTokens(accessTokenExpired, refreshToken string, usernameDB, refreshTokenBDD, refreshTokenExpBDD string) (bool, int, error) {
// 	if accessTokenExpired == "" || refreshToken == "" {
// 		return false, http.StatusUnauthorized, errors.New("token or refresh token is missing")
// 	}

// 	claims, err := paseto.Paseto().ValidateTokenWithoutExpirationCheck(accessTokenExpired)
// 	if err != nil {
// 		return false, http.StatusUnauthorized, err
// 	}

// 	if claims.Username != usernameDB {
// 		return false, http.StatusUnauthorized, errors.New("username mismatch")
// 	}

// 	status, err := refresh.ValidateRefreshToken(refreshToken, refreshTokenExpBDD)
// 	if err != nil {
// 		return false, status, err
// 	}

// 	return true, http.StatusOK, nil
// }
