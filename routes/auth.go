package routes

import (
	"errors"
	"net/http"

	"guardian/data"
	"guardian/exceptions"
	"guardian/security"

	"github.com/gin-gonic/gin"
)

func getTokens(username string, rights int) (token string, refreshToken string, err error) {
	// Création des tokens
	token, errToken := security.CreateToken(username, rights, false)
	refreshToken, errRefreshToken := security.CreateToken(username, security.NONE, true)

	// Gestion d'erreur
	if errToken != nil || errRefreshToken != nil {
		return "", "", errors.New("failed to create tokens")
	}

	return token, refreshToken, nil
}

/**
* Fonction d'enregistrement d'un utilisateur
*
* @param {*gin.Context} c - Context de la requête
 */
func RegisterRoute(c *gin.Context) {
	// Type du body attendu
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Si on n'arrive pas à caster le body, il est mal formé, on renvoie une erreur
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	hashedPassword := security.HashPassword(input.Password)

	// Enregistrement de l'utilisateur
	success, dataError := data.RegisterQuery(input.Username, hashedPassword, security.USER)
	if !success {
		if dataError.Code == exceptions.SQL_ERROR_DUPLICATE {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": dataError.Message})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": dataError.Message})
		}
		return
	} else {
		// Récupération des tokens
		token, refreshToken, err := getTokens(input.Username, security.USER)
		// Gestion des erreurs
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to create tokens"})
			return
		}
		c.IndentedJSON(http.StatusCreated, gin.H{"token": token, "refreshToken": refreshToken})
		return
	}
}

/**
* Fonction de login de l'utilisateur, réponds avec un token d'accès et un token de refresh
*
* @param {*gin.Context} c - Context de la requête
 */
func LoginRoute(c *gin.Context) {
	// Type du body attendu
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Rights   int    `json:"rigths"`
	}

	// Si on n'arrive pas à caster le body, il est mal formé, on renvoie une erreur
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Récupération de l'utilisateur en BDD
	userBdd, err := data.GetUser(input.Username, input.Password)
	// Gestion d'erreur
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	connected := security.HashMatchesPassword(userBdd.Password, input.Password)

	// Erreur de connexion
	if !connected {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	// Création des tokens
	token, refreshToken, err := getTokens(input.Username, userBdd.Rights)

	// Gestion d'erreur
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token, "refreshToken": refreshToken})
}