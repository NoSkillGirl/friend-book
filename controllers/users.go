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
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	errType, err := models.DeleteUser(ctx, userID)

	if err != nil {
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

// //////////////////////////////////////////////////////////////////////////////////////////////////
// type DeleteUserResponse struct {
// 	Message string `json:"message"`
// }

// func DeleteUser(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	w.WriteHeader(http.StatusOK)
// 	userIDString := vars["user_id"]
// 	userID, err := strconv.Atoi(userIDString)

// 	if err != nil {
// 		// do something
// 	}
// 	resp := DeleteUserResponse{}

// 	serverErr := models.DeleteUser(userID)
// 	if serverErr == true {
// 		resp.Message = "Unsuccessful. Internal server Error"
// 		json.NewEncoder(w).Encode(resp)
// 		return
// 	}

// 	resp.Message = "User deleted successfully"
// 	json.NewEncoder(w).Encode(resp)
// }

type ShowUserResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	EmailID string `json:"email"`
	PhoneNo string `json:"phone_no"`
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
		resp.Error.Message = "Unable to parse userID"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	user, errType, err := models.GetUser(ctx, userID)

	if err != nil {
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
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
	return
}

// //////////////////////////////////////////////////////////////////////////////////////////////////
// type FriendRequestRequest struct {
// 	FriendEmailID string `json:"message"`
// }
// type FriendRequestResponse struct {
// 	Message string `json:"message"`
// }

// //sending friend request
// func FriendRequest(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	w.WriteHeader(http.StatusOK)
// 	userIDString := vars["user_id"]
// 	userID, err := strconv.Atoi(userIDString)

// 	if err != nil {
// 		// do something
// 	}

// 	var reqJSON FriendRequestRequest
// 	resp := FriendRequestResponse{}
// 	w.Header().Set("Content-Type", "application/json")

// 	err = json.NewDecoder(r.Body).Decode(&reqJSON)
// 	serverErr := models.FriendRequest(userID, reqJSON.FriendEmailID)

// 	if serverErr == true {
// 		resp.Message = "Unsuccessful. Internal server Error"
// 		json.NewEncoder(w).Encode(resp)
// 		return
// 	}

// 	resp.Message = "friend request sent successfully"
// 	json.NewEncoder(w).Encode(resp)
// }

// ////////////////////////////////////////////////////////////////////////////////////////////////////
// type SearchUser struct {
// 	ID      int    `json:"id"`
// 	Name    string `json:"name"`
// 	EmailID string `json:"email_id"`
// 	PhoneNo string `json:"phone_no"`
// }

// type SearchResponse struct {
// 	Message string       `json:"message"`
// 	Users   []SearchUser `json:"users"`
// }

// func Search(w http.ResponseWriter, r *http.Request) {
// 	// vars := mux.Vars(r)
// 	// w.WriteHeader(http.StatusOK)
// 	// userID := vars["user_id"]
// 	resp := SearchResponse{}

// 	u, serverErr := models.Search()
// 	if serverErr == true {
// 		resp.Message = "Unsuccessful. Internal Server Error"
// 	} else {
// 		resp.Message = "Successful"
// 	}
// 	resp.Friends = f
// 	json.NewEncoder(w).Encode(resp)
// }

/////////////////////////////////////////////////////////////////////////////////////////////////////
type FriendRequestsID struct {
	ID int `json:"id"`
}

type FriendRequestsResponse struct {
	Message          string             `json:"message"`
	FriendRequestIDs []FriendRequestsID `json:"user_ids"`
}

// func FriendRequests(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	w.WriteHeader(http.StatusOK)
// 	userIDString := vars["user_id"]
// 	userID, err := strconv.Atoi(userIDString)

// 	if err != nil {
// 		// do something
// 	}
// 	resp := FriendRequestsResponse{}

// 	fr, serverErr := models.FriendRequests(userID)
// 	if serverErr == true {
// 		resp.Message = "Unsuccessful. Internal Server Error"
// 	} else if len(fr) > 0 {
// 		resp.Message = "Successful"
// 	} else {
// 		resp.Message = "Successful. No friend request found"
// 	}

// 	resp.FriendRequestIDs = fr
// 	json.NewEncoder(w).Encode(resp)
// }

/////////////////////////////////////////////////////////////////////////////////////////////////////////
