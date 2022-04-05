package configs

type JwtKeyConfig struct {
	PublicKeyPath  string `json:"public_key_path" env:"JWT_PUBLIC"`
	PrivateKeyPath string `json:"private_key_path" env:"JWT_PRIVATE"`
}

func GetJwtKeys() (privateKey string, publicKey string) {
	return app.Jwt.PrivateKeyPath, app.Jwt.PublicKeyPath
}