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
	// Variable de multiplication des heures initialisée à 1
	multiplicateur := 1
	// Si en refresh, on fait * 24 heures et * 30 jours
	if refresh {
		multiplicateur = 30 * 24
	}
	// Si pas en refresh, on injecte les droits
	if refresh {
		rights = NONE
	}
	// Création du token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(multiplicateur)).Unix(),
		"rights":   rights,
	})
	// Initialisation de la variable contenant la clé d'accès à la variable d'environnement
	envVariable := "GO_SECRET_KEY_ACCESS_TOKEN"
	// Si en mode refresh, on change la clé de variable
	if refresh {
		envVariable = "GO_SECRET_KEY_REFRESH_TOKEN"
	}
	// Attribution du token et chiffrement avec la variable d'environnement correspondante
	tokenString, err := token.SignedString([]byte(os.Getenv(envVariable)))
	// Si un erreur, on renvoie une string vide et l'objet erreur
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Sinon, on renvoie le token et null
	return &tokenString, nil
}
