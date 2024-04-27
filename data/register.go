package data

import (
	"github.com/go-sql-driver/mysql"

	"guardian/exceptions"
)

/**
* Crée un utilisateur dans la base de données
*
* @param {string} username - Nom d'utilisateur
* @param {string} password - Mot de passe
* @param {int} rights - Niveau de droits
*
* @return {DataPackageError} err - Erreur renvoyée par la fonction
 */
func RegisterQuery(username string, password string, rights int8) (err *exceptions.DataPackageError) {
	tx, errTx := db.Begin()
	if errTx != nil {
		return &exceptions.DataPackageError{Message: "Unable to start transaction", Code: exceptions.SQL_ERROR_LAMBDA}
	}

	_, errEx := tx.Exec(`INSERT INTO auth VALUES (?, ?, ?)`, username, password, rights)
	if errEx != nil {
		if mysqlErr, ok := errEx.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				return &exceptions.DataPackageError{Message: "Duplicate key insertion", Code: exceptions.SQL_ERROR_DUPLICATE}
			} else {
				return &exceptions.DataPackageError{Message: "SQL error", Code: exceptions.SQL_ERROR_LAMBDA}
			}
		} else {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return &exceptions.DataPackageError{Message: "Unable to rollback", Code: exceptions.SQL_ERROR_TRANS_ROLLBACK}
			}
			return &exceptions.DataPackageError{Message: "Internal error", Code: exceptions.ERROR_LAMBDA}
		}
	}

	errTx = tx.Commit()
	if errTx != nil {
		return &exceptions.DataPackageError{Message: "Unable to commit transaction", Code: exceptions.ERROR_LAMBDA}
	}

	return nil
}
