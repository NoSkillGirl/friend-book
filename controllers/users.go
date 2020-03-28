package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/NoSkillGirl/friend-book/constants"
	"github.com/NoSkillGirl/friend-book/models"
	"github.com/gorilla/mux"
)

// ///////////////////////////////////////////////////////////////////////////////////////////////////
type UpdateUserRequest struct {
	Name     string `json:"name"`
	PhoneNo  string `json:"phone"`
	Password string `json:"password"`
}

type UpdateUserResponse struct {
	Message string `json:"message"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Req Obj
	var reqJSON UpdateUserRequest

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	splitReqPath := strings.Split(r.URL.Path, "/")

	lastReqString := splitReqPath[len(splitReqPath)-2]

	var userIDString string
	if lastReqString == "me" {
		userIDString = r.Header.Get("user_id")
	} else {
		vars := mux.Vars(r)
		userIDString = vars["user_id"]
	}

	if userIDString != r.Header.Get("user_id") {
		resp.Error.Message = "You cannot update other user details"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if userIDString == "" {
		resp.Error.Message = "userID missing"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		log.Println("UpdateUser", err)
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Req Decode
	err = json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		resp.Error.Message = "Unable to Parse Request Body"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking if atleast one input is not empty
	if !(reqJSON.Name != "" || reqJSON.PhoneNo != "" || reqJSON.Password != "") {
		resp.Error.Message = "One ore more required fields are missing."
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var passwordSalt string
	if reqJSON.Password != "" {
		passwordSalt, err = hashPassword(reqJSON.Password)
		if err != nil {
			log.Println(err)
			resp.Error.Message = "Unable to create a salt of your password"
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	errType, err := models.UpdateUser(ctx, userID, reqJSON.Name, reqJSON.PhoneNo, passwordSalt)

	if err != nil {
		log.Println(err)
		if errType == constants.ErrorDatabaseUpdateZeroRowsAffected {
			resp.Error.Message = "Looks like database is already up to date"
			resp.Error.Type = constants.ErrorDatabaseUpdateZeroRowsAffected
			resp.Error.Code = http.StatusUnprocessableEntity
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "Error Occured while updating the data"
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	// Res Obj
	successResponse := SignupResponse{}
	successResponse.Success = true
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	splitReqPath := strings.Split(r.URL.Path, "/")

	lastReqString := splitReqPath[len(splitReqPath)-2]

	var userIDString string
	if lastReqString == "me" {
		userIDString = r.Header.Get("user_id")
	} else {
		vars := mux.Vars(r)
		userIDString = vars["user_id"]
	}

	if userIDString != r.Header.Get("user_id") {
		resp.Error.Message = "You cannot delete other users"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if userIDString == "" {
		resp.Error.Message = "userID missing"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		log.Println("DeleteUser", err)
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	errType, err := models.DeleteUser(ctx, userID)

	if err != nil {
		log.Println(err)
		log.Println("Error Occured while deleting the user", errType, err)
		resp.Error.Message = "Error Occured while deleting the user"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Res Obj
	successResponse := SignupResponse{}
	successResponse.Success = true
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}

type ShowUserResponse struct {
	ID      int             `json:"id"`
	Name    string          `json:"name"`
	EmailID string          `json:"email"`
	PhoneNo string          `json:"phone_no"`
	Friends []models.Friend `json:"friend"`
}

// ///////////////////////////////////////////////////////////////////////////////////////////////////
func ShowUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	splitReqPath := strings.Split(r.URL.Path, "/")

	lastReqString := splitReqPath[len(splitReqPath)-1]

	var userIDString string
	if lastReqString == "me" {
		userIDString = r.Header.Get("user_id")
	} else {
		vars := mux.Vars(r)
		userIDString = vars["user_id"]
	}

	if userIDString == "" {
		resp.Error.Message = "userID missing"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		log.Println("ShowUser", err)
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, errType, err := models.GetUser(ctx, userID)

	if err != nil {
		log.Println(err)
		if errType == constants.ErrorDatabaseUserNotFound {
			resp.Error.Message = "This user doesn't exist in our database, please check your user id!"
			resp.Error.Type = constants.ErrorDatabaseUserNotFound
			resp.Error.Code = http.StatusUnprocessableEntity
			log.Println(resp, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "Unable to fetch this User"
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	successResponse := ShowUserResponse{
		ID:      user.ID,
		Name:    user.Name,
		EmailID: user.EmailID,
		PhoneNo: user.PhoneNo,
		Friends: user.Friends,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
	return
}

// //////////////////////////////////////////////////////////////////////////////////////////////////
type FriendRequestRequest struct {
	FriendEmailID string `json:"email"`
}

//sending friend request
func FriendRequest(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	userIDString := r.Header.Get("user_id")

	if userIDString == "" {
		resp.Error.Message = "userID missing"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		log.Println("FriendRequest", err)
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var reqJSON FriendRequestRequest

	w.Header().Set("Content-Type", "application/json")

	// Req Decode
	err = json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		resp.Error.Message = "Unable to Parse Request Body"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking all the inputs
	if reqJSON.FriendEmailID == "" {
		resp.Error.Message = "One ore more required fields are missing."
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validating email-id
	emailValidation := validEmailID(reqJSON.FriendEmailID)
	if emailValidation == false {
		resp.Error.Message = "Invalid email-id"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// err = json.NewDecoder(r.Body).Decode(&reqJSON)
	errType, err := models.FriendRequest(ctx, userID, reqJSON.FriendEmailID)

	if err != nil {
		log.Println(err)
		if errType == constants.ErrorDatabaseUserNotFound {
			resp.Error.Message = "Your friend doesn't exist in our database!"
			resp.Error.Type = constants.ErrorDatabaseUserNotFound
			resp.Error.Code = http.StatusUnprocessableEntity
			log.Println(resp, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		} else if errType == constants.ErrorDatabaseDuplicate {
			resp.Error.Message = "Looks like you have already sent the request to this user."
			resp.Error.Type = constants.ErrorDatabaseDuplicate
			resp.Error.Code = http.StatusConflict
			log.Println(resp, err)
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "Unable to send friend req"
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	// Res Obj
	successResponse := SignupResponse{}
	successResponse.Success = true
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}

// ////////////////////////////////////////////////////////////////////////////////////////////////////
type SearchUser struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	EmailID string `json:"email_id"`
	PhoneNo string `json:"phone_no"`
}

type SearchResponse struct {
	Data []models.Friend `json:"data"`
}

func Search(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	var name, email, phoneNo string
	nameArr, ok := r.URL.Query()["name"]
	if !(!ok || len(nameArr[0]) < 1) {
		name = nameArr[0]
	}

	emailArr, ok := r.URL.Query()["email"]
	if !(!ok || len(emailArr[0]) < 1) {
		email = emailArr[0]
	}

	phoneNoArr, ok := r.URL.Query()["phone_no"]
	if !(!ok || len(phoneNoArr[0]) < 1) {
		phoneNo = phoneNoArr[0]
	}

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	data, errType, err := models.Search(ctx, name, email, phoneNo)

	if err != nil {
		log.Println("Search - ", data, errType, err)
		resp.Error.Message = "Unable to search"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return

	}

	successResp := SearchResponse{
		Data: data,
	}

	json.NewEncoder(w).Encode(successResp)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
type FriendRequestsID struct {
	ID int `json:"id"`
}

type FriendRequestsResponse struct {
	Message string                          `json:"message"`
	Data    []models.FriendRequestsResponse `json:"data"`
}

func FriendRequests(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	userIDString := r.Header.Get("user_id")

	if userIDString == "" {
		resp.Error.Message = "userID missing"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		log.Println("FriendRequests", err)
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	data, errType, err := models.FriendRequests(ctx, userID)

	if err != nil {
		log.Println(err)
		log.Println("FriendRequests - ", data, errType, err)
		if errType == constants.ErrorDatabaseUserNotFound {
			resp.Error.Message = "Your friend doesn't exist in our database!"
			resp.Error.Type = constants.ErrorDatabaseUserNotFound
			resp.Error.Code = http.StatusUnprocessableEntity
			log.Println(resp, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "Unable to send friend req"
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	successResponse := FriendRequestsResponse{}
	if len(data) > 0 {
		successResponse.Message = "Successful"
	} else {
		successResponse.Message = "Successful. No friend request found"
	}

	successResponse.Data = data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
type ActOnFriendRequestRequest struct {
	EmailIDs []string `json:"email_ids"`
	Action   string   `json:"action"`
}

type ActOnFriendRequestResponse struct {
	Data []models.ActOnFriendRequestResponse `json:"data"`
}

func ActOnFriendRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	userIDString := r.Header.Get("user_id")

	if userIDString == "" {
		resp.Error.Message = "userID missing"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		log.Println("ActOnFriendRequest", err)
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var reqJSON ActOnFriendRequestRequest

	// Req Decode
	err = json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		resp.Error.Message = "Unable to Parse Request Body"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking all the inputs
	if len(reqJSON.EmailIDs) == 0 || reqJSON.Action == "" {
		resp.Error.Message = "One ore more required fields are missing."
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validating email-ids
	for _, emailID := range reqJSON.EmailIDs {
		emailValidation := validEmailID(emailID)
		if emailValidation == false {
			resp.Error.Message = "Invalid email-id"
			resp.Error.Type = constants.ErrorValidation
			resp.Error.Code = http.StatusUnprocessableEntity
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	// validation for action
	if !(reqJSON.Action == "accept" || reqJSON.Action == "reject") {
		resp.Error.Message = "Invalid action"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	data, errType, err := models.ActOnFriendRequest(ctx, userID, reqJSON.EmailIDs, reqJSON.Action)
	log.Println(data, errType, err)

	if err != nil {
		log.Println(err)
		if errType == constants.ErrorDatabaseUpdateZeroRowsAffected {
			resp.Error.Message = "All these friend requests are already accepted/rejected"
			resp.Error.Type = constants.ErrorDatabaseUserNotFound
			resp.Error.Code = http.StatusUnprocessableEntity
			log.Println(resp, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "Unable to send friend req"
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	// Res Obj
	successResponse := ActOnFriendRequestResponse{}
	successResponse.Data = data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}

type RemoveFriendRequest struct {
	EmailID string `json:"email"`
}

func RemoveFriend(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	userIDString := r.Header.Get("user_id")

	if userIDString == "" {
		resp.Error.Message = "userID missing"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userID, err := strconv.Atoi(userIDString)

	if err != nil {
		log.Println("RemoveFriend", err)
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	var reqJSON RemoveFriendRequest

	// Req Decode
	err = json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		resp.Error.Message = "Unable to Parse Request Body"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking all the inputs
	if reqJSON.EmailID == "" {
		resp.Error.Message = "One ore more required fields are missing."
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validating email-ids

	emailValidation := validEmailID(reqJSON.EmailID)
	if emailValidation == false {
		resp.Error.Message = "Invalid email-id"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	errType, err := models.RemoveFriend(ctx, userID, reqJSON.EmailID)
	log.Println(errType, err)

	if err != nil {
		log.Println(err)
		if errType == constants.ErrorDatabaseUserNotFound {
			resp.Error.Message = "You don't have a friend with this emailID"
			resp.Error.Type = constants.ErrorDatabaseUserNotFound
			resp.Error.Code = http.StatusUnprocessableEntity
			log.Println(resp, err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "Unable to send friend req"
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	// Res Obj
	successResponse := SignupResponse{
		Success: true,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)

}
