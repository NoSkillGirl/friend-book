package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	//mysql

	"github.com/NoSkillGirl/friend-book/constants"
	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx context.Context
	db  *sql.DB
)

const mySQLHost = "localhost"

var mySQLConnection = fmt.Sprintf("root:@tcp(%s)/friend_book", mySQLHost)

////////////////////////////////////////////////////////////////////////////////////////////////////////
type User struct {
	ID       int
	Name     string
	EmailID  string
	PhoneNo  string
	Password string
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
func DeleteUser(ctx context.Context, userID int) (errType string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		return constants.ErrorDatabaseConnection, err
	}
	defer db.Close()

	delForm, err := db.Prepare("delete from users where id=?")
	_, err = delForm.Exec(userID)
	if err != nil {
		return constants.ErrorDatabaseDelete, err
	}

	return "", nil
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
func FriendRequest(userID int, friendemailID string) (serverErr bool) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer db.Close()

	//get friends id from database
	var friendID int
	//
	//
	//
	//
	//

	//insert the data in database
	insert, err := db.Query(
		`insert into friend_requests (requestor_is, friends_id, status) VALUES (?, ?, "pending")`,
		userID, friendID,
	)

	if err != nil {
		fmt.Println("Error occured while inserting user details in the database", err)
		return true
	}
	defer insert.Close()

	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
type Friend struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	EmailID string `json:"email_id"`
	PhoneNo string `json:"phone_no"`
}

func Search() (users []Friend, serverErr bool) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		fmt.Println(err)
		return users, true
	}
	defer db.Close()

	search, err := db.Query(`select * from users`)
	if err != nil {
		fmt.Println(err)
		return users, true
	}
	defer search.Close()

	for search.Next() {
		u := Friend{}
		err = search.Scan(&u.ID, &u.Name, &u.EmailID, &u.PhoneNo)

		if err != nil {
			fmt.Println(err)
			return users, true
		}
		users = append(users, u)
	}
	return users, false
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////
// type FriendRequestsID struct {
// 	ID int `json:"id"`
// }

// func FriendRequests(userID int) (requestIDs []FriendRequestsID, serverErr bool) {
// 	db, err := sql.Open("mysql", mySQLConnection)
// 	if err != nil {
// 		fmt.Println(err)
// 		return requestIDs, true
// 	}
// 	defer db.Close()

// 	search, err := db.Query(
// 		`select friends_id from friend_requests where requestor_is = ? and status = "active" or status = "pending"`,
// 		userID,
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return requestIDs, true
// 	}
// 	defer search.Close()

// 	for search.Next() {
// 		u := Friend{}
// 		err = search.Scan(&u.ID)

// 		if err != nil {
// 			fmt.Println(err)
// 			return requestIDs, true
// 		}
// 		requestIDs = append(requestIDs, u)
// 	}
// 	return requestIDs, false
// }

////////////////////////////////////////////////////////////////////////////////////////////////////////

func GetUser(ctx context.Context, userID int) (user User, errType string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		return user, constants.ErrorDatabaseConnection, err
	}
	defer db.Close()

	err = db.QueryRow("select name, email_id, phone_no from users where id = ?", userID).Scan(&user.Name, &user.EmailID, &user.PhoneNo)

	if err != nil {
		return user, constants.ErrorDatabaseUserNotFound, err
	}

	user.ID = userID
	return user, "", err
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
func UpdateUser(ctx context.Context, userID int, name string, phoneNo string, password string) (errType string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		return constants.ErrorDatabaseConnection, err
	}
	defer db.Close()

	queryBuilder := "update users set "

	if name != "" {
		queryBuilder += fmt.Sprintf("name='%s'", name)
	}

	if name != "" && phoneNo != "" {
		queryBuilder += fmt.Sprintf(",phone_no='%s'", phoneNo)
	} else if phoneNo != "" {
		queryBuilder += fmt.Sprintf("phone_no='%s'", phoneNo)
	}

	if (name != "" || phoneNo != "") && password != "" {
		queryBuilder += fmt.Sprintf(",password='%s'", password)
	} else if password != "" {
		queryBuilder += fmt.Sprintf("password='%s'", password)
	}

	queryBuilder += fmt.Sprintf(" where id = %v", userID)

	result, err := db.ExecContext(ctx, queryBuilder)
	if err != nil {
		return constants.ErrorDatabaseUpdate, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return constants.ErrorDatabaseUpdate, err
	}
	if rows != 1 {
		return constants.ErrorDatabaseUpdateZeroRowsAffected, errors.New("Zero rows affected")
	}

	return "", nil
}
