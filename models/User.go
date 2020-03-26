package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

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
func FriendRequest(ctx context.Context, userID int, friendemailID string) (errType string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		return constants.ErrorDatabaseConnection, err
	}
	defer db.Close()

	//get friends id from database
	var friendID int

	err = db.QueryRow("select id from users where email_id = ?", friendemailID).Scan(&friendID)

	if err != nil {
		return constants.ErrorDatabaseUserNotFound, err
	}

	//insert the data in database
	insert, err := db.Query(
		`insert into friend_requests (requestor_id, friend_id, status) VALUES (?, ?, "pending")`,
		userID, friendID,
	)

	if err != nil {
		log.Println("Error occured while inserting friend request details in the database", err)
		return constants.ErrorDatabaseInsert, err
	}
	defer insert.Close()

	return "", nil
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
type FriendRequestsResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func FriendRequests(ctx context.Context, userID int) (requestIDs []FriendRequestsResponse, errType string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		return requestIDs, constants.ErrorDatabaseConnection, err
	}
	defer db.Close()

	search, err := db.Query(
		`select friend_id from friend_requests where requestor_id = ? and status = "pending"`,
		userID,
	)

	if err != nil {
		return requestIDs, constants.ErrorDatabaseSelect, err
	}
	defer search.Close()

	friendIDs := make([]string, 0)

	for search.Next() {
		var ID string
		err = search.Scan(&ID)

		if err != nil {
			return requestIDs, constants.ErrorDatabaseSelect, err
		}
		friendIDs = append(friendIDs, ID)
	}

	queryString := fmt.Sprintf("select id, email_id from users where id in (%s)", "'"+strings.Join(friendIDs[:], "','")+"'")

	search, err = db.Query(queryString)

	if err != nil {
		return requestIDs, constants.ErrorDatabaseSelect, err
	}
	defer search.Close()

	for search.Next() {
		var frr FriendRequestsResponse
		err = search.Scan(&frr.ID, &frr.Email)

		if err != nil {
			return requestIDs, constants.ErrorDatabaseSelect, err
		}
		requestIDs = append(requestIDs, frr)
	}

	return requestIDs, "", nil
}

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

type ActOnFriendRequestResponse struct {
	Email  string `json:"email"`
	Status string `json:"status"`
}

///////////////////////////////////////////////////////////////////////////////////////////////////////
func ActOnFriendRequest(ctx context.Context, userID int, emailIDs []string, action string) (data []ActOnFriendRequestResponse, errType string, err error) {
	db, err := sql.Open("mysql", mySQLConnection)
	if err != nil {
		return data, constants.ErrorDatabaseConnection, err
	}
	defer db.Close()

	var status string
	if action == "accept" {
		status = "accepted"
	} else {
		status = "rejected"
	}

	// get friend ids
	queryString := fmt.Sprintf("select id, email_id from users where email_id in (%s)", "'"+strings.Join(emailIDs[:], "','")+"'")

	search, err := db.Query(queryString)

	if err != nil {
		return data, constants.ErrorDatabaseSelect, err
	}
	defer search.Close()

	var friendIDs []string
	var friendEmails []string
	for search.Next() {
		var friendID string
		var friendEmail string
		err = search.Scan(&friendID, &friendEmail)

		if err != nil {
			return data, constants.ErrorDatabaseSelect, err
		}
		friendIDs = append(friendIDs, friendID)
		friendEmails = append(friendEmails, friendEmail)
	}

	// check if any friend requests are pending
	queryString = fmt.Sprintf("select count(*) from friend_requests where status = 'pending' and requestor_id = '%v' and friend_id in ('%s')", userID, strings.Join(friendIDs[:], "','"))

	var count int
	err = db.QueryRow(queryString).Scan(&count)

	if err != nil {
		return data, constants.ErrorDatabaseUpdate, err
	}

	if count == 0 {
		return data, constants.ErrorDatabaseUpdateZeroRowsAffected, errors.New("Zero rows to update")
	}

	// update friend_requests
	query := fmt.Sprintf("update friend_requests set status = '%s' where requestor_id = '%v' and friend_id in ('%s')", status, userID, strings.Join(friendIDs[:], "','"))

	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return data, constants.ErrorDatabaseUpdate, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return data, constants.ErrorDatabaseUpdate, err
	}
	if rows == 0 {
		return data, constants.ErrorDatabaseUpdateZeroRowsAffected, errors.New("Zero rows affected")
	}

	for _, emailID := range friendEmails {
		var afr ActOnFriendRequestResponse
		afr.Email = emailID
		afr.Status = status
		data = append(data, afr)
	}

	return data, "", nil
}
