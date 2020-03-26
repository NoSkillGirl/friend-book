package models

import (
	"context"
	"database/sql"
	"errors"
	"log"

	//mysql
	"github.com/NoSkillGirl/friend-book/constants"
	_ "github.com/go-sql-driver/mysql"
)

func Login(ctx context.Context, emailID, password string) (userID string, passwordSalt string, errType string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		return userID, passwordSalt, constants.ErrorDatabaseConnection, err
	}
	defer db.Close()

	err = db.QueryRow("select id, password from users where email_id = ?", emailID).Scan(&userID, &passwordSalt)

	if err != nil {
		return userID, passwordSalt, constants.ErrorDatabaseEmailNotFound, err
	}

	return userID, passwordSalt, "", err
}

//SignUp function
func SignUp(ctx context.Context, name, emailID, phoneNo, password string) (typeOfError string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	// if there is an error opening the connection, handle it
	if err != nil {
		return constants.ErrorDatabaseConnection, err
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	var count int
	err = db.QueryRowContext(ctx, "select count(*) from users where email_id=? or phone_no=?", emailID, phoneNo).Scan(&count)

	if err != nil {
		return constants.ErrorDatabaseSelect, err
	}

	if count != 0 {
		return constants.ErrorDatabaseDuplicate, errors.New("duplicate error")
	}

	insert, err := db.Query(
		`insert into users (name, email_id, phone_no, password) VALUES (?, ?, ?, ?)`,
		name, emailID, phoneNo, password,
	)

	if err != nil {
		log.Println("Error occured while inserting user details in the database", err)
		return constants.ErrorDatabaseInsert, err
	}
	defer insert.Close()
	return "", nil
}
