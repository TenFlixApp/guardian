package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"guardian/data"
	"guardian/exceptions"
	"guardian/security"

	"github.com/gin-gonic/gin"
)

type Input struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Rights   int    `json:"rights"`
}

/**
* Création des tokens
*
* @param {string} username - Nom d'utilisateur
* @param {int} rights - Niveau des droits
*
* @returns {string} token - Token d'accès
* @returns {string} refreshToken - Token de refresh
* @returns {error} err - Erreur lors de la fonction
 */
func getTokens(username string, rights int) (token *string, refreshToken *string, err error) {
	token, errToken := security.CreateToken(username, rights, false)
	refreshToken, errRefreshToken := security.CreateToken(username, security.NONE, true)

	if errToken != nil || errRefreshToken != nil {
		return nil, nil, errors.New("failed to create tokens")
	}

	return token, refreshToken, nil
}

/**
* Fonction d'enregistrement d'un utilisateur
*
* @param {*gin.Context} c - Context de la requête
 */
func RegisterRoute(c *gin.Context) {
	var input Input

	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword := security.HashPassword(input.Password)

	dataError := data.RegisterQuery(input.Username, hashedPassword, security.USER)
	if dataError != nil {
		if dataError.Code == exceptions.SQL_ERROR_DUPLICATE {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": dataError.Message})
		} else {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": dataError.Message})
		}
		return
	} else {
		token, refreshToken, err := getTokens(input.Username, security.USER)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Failed to create tokens"})
			return
		}

		metrics := map[string]interface{}{
			"username":      input.Username,
			"register_date": time.Now().Format(time.RFC3339),
		}
		jsonMetrics, err := json.Marshal(metrics)
		if err == nil {
			_, err = http.Post(os.Getenv("COLLECTOR_ROUTE")+"metrics/register", "application/json", bytes.NewBuffer(jsonMetrics))
			if err != nil {
				log.Println("Failed to push register metrics", err)
				return
			}
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
	var input Input
	if err := c.BindJSON(&input); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userBdd, err := data.GetUser(input.Username, input.Password)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid user"})
		return
	}

	connected := security.HashMatchesPassword(userBdd.Password, input.Password)

	if !connected {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	token, refreshToken, err := getTokens(input.Username, userBdd.Rights)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tokens"})
	}

	metrics := map[string]interface{}{
		"username":   input.Username,
		"login_date": time.Now().Format(time.RFC3339),
	}
	jsonMetrics, err := json.Marshal(metrics)
	if err == nil {
		_, err = http.Post(os.Getenv("COLLECTOR_ROUTE")+"metrics/login", "application/json", bytes.NewBuffer(jsonMetrics))
		if err != nil {
			log.Println("Failed to push login metrics", err)
			return
		}
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token, "refreshToken": refreshToken})
}
