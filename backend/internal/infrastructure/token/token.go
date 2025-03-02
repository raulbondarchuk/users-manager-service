package token

// TODO: Realizar una func (expiredToken, refresh) para generar un nuevo par de tokens.
// Hay que encontrar usuario por refresh token, comprobar si ha expirado, generar un nuevo par de tokens y actualizar el usuario con los nuevos tokens.

func GeneratePairTokens(expiredAccessToken string, refreshToken string) (string, string, error) {
	return "", "", nil
}
