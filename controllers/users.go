package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/NoSkillGirl/friend-book/models"
	"github.com/gorilla/mux"
)

// HealthCheck - health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}

type UserSignupRequest struct {
	Name     string `json:"name"`
	EmailID  string `json:"email_id"`
	PhoneNo  string `json:"phone_no"`
	Password string `json:"password"`
}

type UserSignupResponse struct {
	Message string `json:"message"`
}

func UserSignup(w http.ResponseWriter, r *http.Request) {
	// Req Obj
	var reqJSON UserSignupRequest
	// Res Obj
	resp := UserSignupResponse{}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.Body)

	// Req Decode
	err := json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		fmt.Println(err)
		resp.Message = "unsuccessful"
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking all the inputs
	if reqJSON.Name == "" || reqJSON.EmailID == "" || reqJSON.PhoneNo == "" || reqJSON.Password == "" {
		fmt.Println("One ore more detail/s is/ are missing")
		resp.Message = "Unsuccessful. One ore more detail/s is/are missing"
		json.NewEncoder(w).Encode(resp)
		return
	}
	//validating email-id
	emailValidation := validEmailID(reqJSON.EmailID)
	if emailValidation == false {
		fmt.Println("Invalid email-id")
		resp.Message = "Unsuccessful. Invalid email-id"
		json.NewEncoder(w).Encode(resp)
		return
	}
	//validating phone number
	phoneValidation := validPhoneNo(reqJSON.PhoneNo)
	if phoneValidation == false {
		fmt.Println("Invalid phone number")
		resp.Message = "Unsuccessful. Invalid phone number"
		json.NewEncoder(w).Encode(resp)
		return
	}
	insertErr, duplicateErr := models.UserSignup(reqJSON.Name, reqJSON.EmailID, reqJSON.PhoneNo, reqJSON.Password)
	if insertErr == true {
		resp.Message = "Unsuccessful. Internal server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	if duplicateErr == true {
		fmt.Println("Duplicate found")
		resp.Message = "Unsuccessful. Duplicate found"
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp.Message = "User sucessfully created"
	json.NewEncoder(w).Encode(resp)
}

func validEmailID(emailID string) bool {
	r := regexp.MustCompile("^([\w-\.]+)@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.)|(([\w-]+\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\]?)$")
	return r.MatchString(emailID)
}

func validPhoneNo(phoneNo string) bool {
	r := regexp.MustCompile("^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$")
	return r.MatchString(phoneNo)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
type UserLoginRequest struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Message string `json:"message"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	// Req Obj
	var reqJSON UserLoginRequest
	// Res Obj
	resp := UserLoginResponse{}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.Body)

	// Req Decode
	err := json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		fmt.Println(err)
		resp.Message = "Login unsuccessful"
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking the inputs
	if reqJSON.EmailID == "" || reqJSON.Password == "" {
		fmt.Println("Not enough details")
		resp.Message = "Unsuccessful. Not enough details"
		json.NewEncoder(w).Encode(resp)
		return
	}
	//validating email-id
	emailValidation := validEmailID(reqJSON.EmailID)
	if emailValidation == false {
		fmt.Println("Invalid email-id")
		resp.Message = "Unsuccessful. Invalid email-id"
		json.NewEncoder(w).Encode(resp)
		return
	}

	u, serverErr := models.UserLogin(reqJSON.EmailID, reqJSON.Password)

	if serverErr == true {
		resp.Message = "Unsuccessful. Internal server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}
	if len(u) > 0 {
		resp.Message = "Successful. User found"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Message = "User not found."
	json.NewEncoder(w).Encode(resp)
	return
}

///////////////////////////////////////////////////////////////////////////////////////////////////
type UpdateUserRequest struct {
	Name     string `json:"name"`
	PhoneNo  string `json:"phone_no"`
	Password string `json:"password"`
}

type UpdateUserResponse struct {
	Message string `json:"message"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	userID := vars["user_id"]

	// Req Obj
	var reqJSON UpdateUserRequest
	// Res Obj
	resp := UpdateUserResponse{}
	w.Header().Set("Content-Type", "application/json")
	fmt.Println(r.Body)

	// Req Decode
	err := json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		fmt.Println(err)
		resp.Message = "Unsuccessful"
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking all the inputs
	if reqJSON.Name == "" || reqJSON.PhoneNo == "" || reqJSON.Password == "" {
		fmt.Println("One ore more detail/s is/ are missing")
		resp.Message = "Unsuccessful. One ore more detail/s is/are missing"
		json.NewEncoder(w).Encode(resp)
		return
	}

	serverErr := models.UpdateUser(userID, reqJSON.Name, reqJSON.PhoneNo, reqJSON.Password)
	if serverErr == true {
		resp.Message = "Unsuccessful. Internal server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Message = "User updated successfully"
	json.NewEncoder(w).Encode(resp)
}

//////////////////////////////////////////////////////////////////////////////////////////////////
type DeleteUserResponse struct {
	Message string `json:"message"`
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	userID := vars["user_id"]
	resp := DeleteUserResponse{}

	serverErr := models.DeleteUser(userID)
	if serverErr == true {
		resp.Message = "Unsuccessful. Internal server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Message = "User deleted successfully"
	json.NewEncoder(w).Encode(resp)
}

///////////////////////////////////////////////////////////////////////////////////////////////////
func ShowUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	userID := vars["user_id"]
	fmt.Fprintf(w, "user id : %v\n", userID)

}

//////////////////////////////////////////////////////////////////////////////////////////////////
type FriendRequestRequest struct {
	FriendEmailID string   `json:"message"`
}
type FriendRequestResponse struct {
	Message string   `json:"message"`
}
//sending friend request
func FriendRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	userID := vars["user_id"]
	
	var reqJSON FriendRequestRequest
	resp := FriendRequestResponse{}
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&reqJSON)
	serverErr := models.FriendRequest(userID, reqJSON.FriendEmailID)

	if serverErr == true {
		resp.Message = "Unsuccessful. Internal server Error"
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Message = "friend request sent successfully"
	json.NewEncoder(w).Encode(resp)
}

////////////////////////////////////////////////////////////////////////////////////////////////////
type SearchUser struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	EmailID string `json:"email_id"`
	PhoneNo string `json:"phone_no"`
}

type SearchResponse struct {
	Message string       `json:"message"`
	Users   []SearchUser `json:"users"`
}

func Search(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// w.WriteHeader(http.StatusOK)
	// userID := vars["user_id"]
	resp := SearchResponse{}

	u, serverErr := models.Search()
	if serverErr == true {
		resp.Message = "Unsuccessful. Internal Server Error"
	} else {
		resp.Message = "Successful"
	}
	resp.Friends = f
	json.NewEncoder(w).Encode(resp)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
type FriendRequestsID struct {
	ID int `json:"id"`
}

type FriendRequestsResponse struct {
	Message          string             `json:"message"`
	FriendRequestIDs []FriendRequestsID `json:"user_ids"`
}

func FriendRequests(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	userID := vars["user_id"]
	resp := FriendRequestsResponse{}

	fr, serverErr := models.FriendRequests(userID)
	if serverErr == true {
		resp.Message = "Unsuccessful. Internal Server Error"
	} else if len(fr) > 0 {
		resp.Message = "Successful"
	} else {
		resp.Message = "Successful. No friend request found"
	}

	resp.FriendRequestIDs = fr
	json.NewEncoder(w).Encode(resp)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
