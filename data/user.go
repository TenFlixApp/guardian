package data

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
func GetUser(username string, password string) (*User, error) {
	utilisateur := &User{}

	err := db.QueryRow(`SELECT username, password, rights FROM auth WHERE username = ?`, username).Scan(&utilisateur.Username, &utilisateur.Password, &utilisateur.Rights)
	if err != nil {
		return nil, err
	} else {
		return utilisateur, nil
	}
}

func CountAdminUsers() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM auth WHERE rights = 15").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CountDisabledUsers() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM auth WHERE rights = 0").Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
