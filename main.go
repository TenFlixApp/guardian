package main

/**
* Imports nécessaires pour l'api
 */
import (
	"log"

	"guardian/data"
	"guardian/routes"

	"github.com/gin-contrib/cors"      // CORS
	"github.com/gin-gonic/gin"         // Framework GIN
	_ "github.com/go-sql-driver/mysql" // Package de tokens JWT
	"github.com/joho/godotenv"         // package de lecture de variable d'environnement
)

func main() {
	// Charchement .env
	err := godotenv.Load()
	// Si erreur, on plente
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	// Tentative de connexion à la base de donnée
	data.ConnectToDB()
	// Programmation de la fermeture de la base de données à la fermeture du programme
	defer data.CloseDB()

	// Création du routeur avec le framework GIN
	router := gin.Default()

	// Création de la configuration par défaut des cors du serveur
	configCors := cors.DefaultConfig()

	// Modification des paramètres
	configCors.AllowAllOrigins, configCors.AllowCredentials = true, true
	configCors.AddAllowHeaders("Authorization")
	configCors.AddAllowHeaders("creditential")

	// Application de la nouvelle configuration
	router.Use(cors.New(configCors))

	// Enregistrement de l'utilisateur
	router.POST("/register", routes.RegisterRoute)

	// Connexion de l'utilisateur
	router.POST("/login", routes.LoginRoute)

	// Lancement du serveur
	router.Run(":8080")
}
