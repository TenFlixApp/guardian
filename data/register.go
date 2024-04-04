package data

import (
	"github.com/go-sql-driver/mysql" // Driver MySQL

	"guardian/exceptions" // Package exception
)

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
