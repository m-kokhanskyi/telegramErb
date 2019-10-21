package models

import (
	"log"
	"time"
)

//User type
type User struct {
	ID          int64
	IDChat      int64
	Login       string
	DataSearch  string
	IsSearching bool
	ConnectedAt time.Time
}

// CreateUser - add new user to DB
func (db *MySQLStorage) CreateUser(u *User) {
	result, err := db.con.Exec("INSERT INTO user(id_telegram, login, is_searching) VALUES(?, ?, ?)", u.IDChat, u.Login, u.IsSearching)
	if err != nil {
		panic(err)
	}
	u.ID, err = result.LastInsertId()
}

// GetAllUsers - find users
func (db *MySQLStorage) GetAllUsers() []User {
	result, err := db.con.Query("SELECT id, id_telegram, login, data_search FROM user WHERE is_searching = 1")
	if err != nil {
		panic(err)
	}
	users := []User{}

	for result.Next() {
		u := User{}
		err := result.Scan(&u.ID, &u.IDChat, &u.Login, &u.DataSearch)
		if err != nil {
			log.Fatal(err)
			continue
		}
		users = append(users, u)
	}

	return users
}

// SetDataSearch - find user by id chat
func (db *MySQLStorage) SetDataSearch(u *User) {
	_, err := db.con.Exec("UPDATE user SET data_search = ?, is_searching = ? WHERE id_telegram = ?", u.DataSearch, u.IsSearching, u.IDChat)

	if err != nil {
		panic(err)
	}
}

func (db *MySQLStorage) SetIsSearching(u *User) {
	_, err := db.con.Exec("UPDATE user SET is_searching = ? WHERE id_telegram = ?", u.IsSearching, u.IDChat)

	if err != nil {
		panic(err)
	}
}
