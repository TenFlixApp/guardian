package data

import (
	"errors"
	"fmt"
)

// Type utilisateur
type User struct {
	Username string
	Password string
	Rights   int
}

/**
* Crée un utilisateur dans la base de données
*
* @param {string} username - Nom d'utilisateur
* @param {string} password - Mot de passe
*
* @return {bool} success - Succès de la fonction
 */
func GetUser(username string, password string) (user *User, err error) {
	// Création d'un user vide
	utilisateur := &User{}

	// Récupération du user BDD
	erreur := db.QueryRow(`SELECT username, password, rights FROM auth WHERE username = ?`, username).Scan(&utilisateur.Username, &utilisateur.Password, &utilisateur.Rights)
	// Gestion d'erreur
	if erreur != nil {
		fmt.Printf("Erreur de récupération du user: %v\n", err)
		return &User{}, errors.New("utilisateur non récupéré")
	} else {
		return utilisateur, nil
	}
}
