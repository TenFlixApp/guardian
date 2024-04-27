package data

import "errors"

type UserPw struct {
	Username string `json:"username"`
	Password string
	Rights   int `json:"rights"`
}

type User struct {
	Username string `json:"username"`
	Rights   int    `json:"rights"`
}

/**
* Crée un utilisateur dans la base de données
*
* @param {string} username - Nom d'utilisateur
* @param {string} password - Mot de passe
*
* @return {bool} success - Succès de la fonction
 */
func GetUser(username string, password string) (*UserPw, error) {
	utilisateur := &UserPw{}

	err := db.QueryRow(`SELECT username, password, rights FROM auth WHERE username = ?`, username).Scan(&utilisateur.Username, &utilisateur.Password, &utilisateur.Rights)
	if err != nil {
		return nil, err
	} else {
		return utilisateur, nil
	}
}

func GetUsers() ([]User, error) {
	var users = make([]User, 0)

	rows, err := db.Query(`SELECT username, rights FROM auth`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Username, &user.Rights)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func UpdateUserRights(username string, rights int) error {
	res, err := db.Exec(`UPDATE auth SET rights = ? WHERE username = ?`, rights, username)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
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
