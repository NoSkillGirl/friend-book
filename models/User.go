package models

import (
	"context"
	"database/sql"
	"fmt"

	//mysql
	_ "github.com/go-sql-driver/mysql"
)

var (
	ctx context.Context
	db  *sql.DB
)

const mySQLHost = "34.93.137.151"

var mySQLConnection = fmt.Sprintf("root:password@tcp(%s)/friend_book", mySQLHost)

//UserSignup function
func UserSignup(name, emailID, phoneNo, password string) (insertError, duplicateUser bool) {
	ctx := context.Background()
	db, err := sql.Open("mysql", mySQLConnection)
	// if there is an error opening the connection, handle it
	if err != nil {
		return true, false
	}

	// defer the close till after the main function has finished executing
	defer db.Close()
	var count int
	err = db.QueryRowContext(ctx, "select count(*) from users where email_id=? or phone_no=?", emailID, phoneNo).Scan(&count)

	if err != nil {
		return true, false
	}

	//Duplicate not found
	if count == 0 {
		insert, err := db.Query(
			`insert into users (name, email_is, phone_no, password) VALUES (?, ?, ?, ?)`,
			name, emailID, phoneNo, password,
		)

		if err != nil {
			fmt.Println("Error occured while inserting user details in the database", err)
			return true, false
		}
		defer insert.Close()
		return false, false
	}

	//error = false, duplicate = true
	return false, true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
type User struct {
	ID       int
	Name     string
	EmailID  string
	PhoneNo  string
	Password string
}

//UserLogin function
func UserLogin(emailID, password string) (user []User, serverErr bool) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		fmt.Println(err)
		return user, true
	}
	defer db.Close()

	search, err := db.Query(
		`select * from users where email_id = ? and password =? limit 1`,
		emailID, password,
	)
	if err != nil {
		fmt.Println(err)
		return user, true
	}
	defer search.Close()

	for search.Next() {
		u := User{}
		err = search.Scan(&u.ID, &u.Name, &u.EmailID, &u.PhoneNo, &u.Password)

		if err != nil {
			fmt.Println(err)
			return user, true
		}
		user = append(user, u)
	}
	return user, false
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
func UpdateUser(userID int, name string, phoneNo string, password string) (serverErr bool) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer db.Close()

	insForm, err := db.Prepare("update users set name=?, phone_no=?, password=? where id=?")
	_, err = insForm.Exec(name, phoneNo, password, userID)
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
func DeleteUser(userID int) (serverErr bool) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer db.Close()

	delForm, err := db.Prepare("delete from users where id=?")
	_, err = delForm.Exec(userID)
	if err != nil {
		fmt.Println(err)
		return true
	}

	return false
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
type FriendRequestsID struct {
	ID int `json:"id"`
}

func FriendRequests(userID int) (requestIDs []FriendRequestsID, serverErr bool) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		fmt.Println(err)
		return requestIDs, true
	}
	defer db.Close()

	search, err := db.Query(
		`select friends_id from friend_requests where requestor_is = ? and status = "active" or status = "pending"`,
		userID,
	)
	if err != nil {
		fmt.Println(err)
		return requestIDs, true
	}
	defer search.Close()

	for search.Next() {
		u := Friend{}
		err = search.Scan(&u.ID)

		if err != nil {
			fmt.Println(err)
			return requestIDs, true
		}
		requestIDs = append(requestIDs, u)
	}
	return requestIDs, false
}

////////////////////////////////////////////////////////////////////////////////////////////////////////
