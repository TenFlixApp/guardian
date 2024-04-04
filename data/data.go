package data

import (
	"database/sql" // Package sql
	"errors"       // Erreurs
	"fmt"          // Affichage console
	"log"          // Logs
	"os"           // OS

	"github.com/go-sql-driver/mysql" // Driver MySQL

	"guardian/exceptions" // Package exception
)

// Variable driver BDD
var db *sql.DB

// Type utilisateur
type User struct {
	Username string
	Password string
	Rights   int
}

/**
* Connexion à la base de données
 */
func ConnectToDB() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_CONN_STRING"))
	if err != nil {
		log.Fatal("Unable to create DB handle", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the DB", err)
	}
	log.Println("Connected to the database")
}

/**
* Fermeture de la connexion BDD
 */
func CloseDB() {
	err := db.Close()
	if err != nil {
		log.Fatalln("Error closing the database connection")
	}
}

/**
* Crée un utilisateur dans la base de données
*
* @param {string} username - Nom d'utilisateur
* @param {string} password - Mot de passe
* @param {int} rights - Niveau de droits
*
* @return {bool} success - Succès de la fonction
 */
func RegisterQuery(username string, password string, rights int8) (success bool, erreur exceptions.DataPackageError) {
	// On ouvre une transaction BDD
	tx, err := db.Begin()
	// Si erreur, on plante
	if err != nil {
		// Gestion d'erreur
		return false, exceptions.DataPackageError{Message: "Unable to start transaction", Code: exceptions.SQL_ERROR_LAMBDA}
	}

	// Insertion de l'utilisateur
	_, err = tx.Exec(`INSERT INTO auth VALUES (?, ?, ?)`, username, password, rights)
	// Gestion erreur
	if err != nil {
		// Vérifier si c'est une erreur MySQL
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			// Vérifier si c'est une erreur de clé dupliquée
			if mysqlErr.Number == 1062 {
				// Retour de l'erreur de duplication de clé
				return false, exceptions.DataPackageError{Message: "Duplicate key insertion", Code: exceptions.SQL_ERROR_DUPLICATE}
			} else {
				// Autre type d'erreur MySQL
				return false, exceptions.DataPackageError{Message: "SQL error", Code: exceptions.SQL_ERROR_LAMBDA}
			}
		} else {
			// Autre type d'erreur
			// Tentative de rollback
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return false, exceptions.DataPackageError{Message: "Unable to rollback", Code: exceptions.SQL_ERROR_TRANS_ROLLBACK}
			}
			return false, exceptions.DataPackageError{Message: "Internal error", Code: exceptions.ERROR_LAMBDA}
		}
	}

	// Commit la transaction
	err = tx.Commit()
	// Gestion erreur
	if err != nil {
		// Autre type d'erreur
		return false, exceptions.DataPackageError{Message: "Unable to commit transaction", Code: exceptions.ERROR_LAMBDA}
	}

	// Retourne succès à true
	return true, exceptions.DataPackageError{}
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
