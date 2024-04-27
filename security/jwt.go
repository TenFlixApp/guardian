package security

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

/**
* Fonction de création du token d'accès / refresh
*
* @param {string} username - Nom de l'utilisateur
* @param {bool} refresh - True si on gènère le token de refresh, false si access
 */
func CreateToken(username string, rights int, refresh bool) (*string, error) {
	multiplicateur := 1
	if refresh {
		multiplicateur = 30 * 24
	}

	if refresh {
		rights = NONE
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(multiplicateur)).Unix(),
		"rights":   rights,
	})

	envVariable := "GO_SECRET_KEY_ACCESS_TOKEN"
	if refresh {
		envVariable = "GO_SECRET_KEY_REFRESH_TOKEN"
	}

	tokenString, err := token.SignedString([]byte(os.Getenv(envVariable)))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &tokenString, nil
}
