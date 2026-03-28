package helper

import "golang.org/x/crypto/bcrypt"

// HashPassword genera el hash de una contraseña
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		12, // costo seguro aumentado
	)
	return string(bytes), err
}

// CheckPassword compara contraseña en texto plano con hash
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)
}

/*
err := security.CheckPassword(passwordIngresada, usuarioDB.UsuPassword)
if err != nil {
	return errors.New("credenciales inválidas")
}
*/
