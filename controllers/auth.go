package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/NoSkillGirl/friend-book/constants"
	"github.com/NoSkillGirl/friend-book/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//hashPassword - func to convert the password to salt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//checkPasswordHash - func to comaper the password and salt
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenericErrorResponse - struct for error
type GenericErrorResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    int    `json:"code"`
}

//ErrorResponse - struct for error
type ErrorResponse struct {
	Error GenericErrorResponse `json:"error"`
}

// SignupRequest - req structure for user signup
type SignupRequest struct {
	Name     string `json:"name"`
	EmailID  string `json:"email"`
	PhoneNo  string `json:"phone_no"`
	Password string `json:"password"`
}

// SignupResponse - resp structure for user signup
type SignupResponse struct {
	Success bool `json:"success"`
}

//Signup - function to handle the user signup
func Signup(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	// Req Obj
	var reqJSON SignupRequest
	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	// Req Decode
	err := json.NewDecoder(r.Body).Decode(&reqJSON)
	if err != nil {
		log.Println("Signup - ", err)
		resp.Error.Message = "Unable to Parse Request Body"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validations
	//Checking all the inputs
	if reqJSON.Name == "" || reqJSON.EmailID == "" || reqJSON.PhoneNo == "" || reqJSON.Password == "" {
		resp.Error.Message = "One ore more required fields are missing."
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validating email-id
	emailValidation := validEmailID(reqJSON.EmailID)
	if emailValidation == false {
		resp.Error.Message = "Invalid email-id"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}
	//validating phone number
	phoneValidation := validPhoneNo(reqJSON.PhoneNo)
	if phoneValidation == false {
		resp.Error.Message = "Invalid phone number"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//converting password to salt
	passwordSalt, err := hashPassword(reqJSON.Password)
	if err != nil {
		log.Println(err)
		resp.Error.Message = "Unable to create a salt of your password"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// calling model
	errType, err := models.SignUp(ctx, reqJSON.Name, reqJSON.EmailID, reqJSON.PhoneNo, passwordSalt)
	if err != nil {
		log.Println(err)
		if errType == constants.ErrorDatabaseDuplicate {
			resp.Error.Message = "Looks like you have already registered. Try logging in"
			resp.Error.Type = constants.ErrorAlreadyRegistered
			resp.Error.Code = http.StatusConflict
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "User signup failed."
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

func GetSignUpPage(w http.ResponseWriter, r *http.Request) {
	err := renderHTML(w, "register")
	if err != nil {
		renderHTML(w, "404")
	}
}

//
//LoginRequest - struct
type LoginRequest struct {
	EmailID  string `json:"email"`
	Password string `json:"password"`
}

//LoginResponse - struct
type LoginResponse struct {
	AuthToken string `json:"auth_token"`
}

//Login - function to handle the user login
func Login(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	// Req Obj
	var reqJSON LoginRequest
	// Res Obj
	resp := ErrorResponse{}
	w.Header().Set("Content-Type", "application/json")

	// Req Decode
	err := json.NewDecoder(r.Body).Decode(&reqJSON)
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

	// validations
	//Checking the inputs
	if reqJSON.EmailID == "" || reqJSON.Password == "" {
		resp.Error.Message = "One ore more required fields are missing."
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//validating email-id
	emailValidation := validEmailID(reqJSON.EmailID)
	if emailValidation == false {
		resp.Error.Message = "Invalid email-id"
		resp.Error.Type = constants.ErrorValidation
		resp.Error.Code = http.StatusUnprocessableEntity
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//calling model
	userID, passwordSalt, errType, err := models.Login(ctx, reqJSON.EmailID, reqJSON.Password)
	if err != nil {
		log.Println(err)
		if errType == constants.ErrorDatabaseEmailNotFound {
			resp.Error.Message = "Your email doesn't exist in our database, please check your email id!"
			resp.Error.Type = constants.ErrorDatabaseEmailNotFound
			resp.Error.Code = http.StatusUnprocessableEntity
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Error.Message = "User login failed."
			resp.Error.Type = constants.ErrorInternalServerError
			resp.Error.Code = http.StatusInternalServerError
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}

	//matching password and salt
	match := checkPasswordHash(reqJSON.Password, passwordSalt)
	if !match {
		resp.Error.Message = "User login failed. Incorrect Password"
		resp.Error.Type = constants.ErrorStatusUnauthorized
		resp.Error.Code = http.StatusUnauthorized
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	//generating JWT
	validToken, err := generateJWT(userID)
	if err != nil {
		log.Println(err)
		resp.Error.Message = "Failed to generate Auth token"
		resp.Error.Type = constants.ErrorInternalServerError
		resp.Error.Code = http.StatusInternalServerError
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	successResponse := LoginResponse{
		AuthToken: string(validToken),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
	return
}

type GetLoginPageTemplateData struct {
	StaticDir string `protobuf:"bytes,1,opt,name=static_dir,json=staticDir,proto3" json:"static_dir,omitempty"`
}

// GetLoginPage - for html page
func GetLoginPage(w http.ResponseWriter, r *http.Request) {
	err := renderHTML(w, "login")
	if err != nil {
		renderHTML(w, "404")
	}
}

// HealthCheck - health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func generateJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	tokenString, err := token.SignedString([]byte(userID))

	if err != nil {
		log.Println(err)
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

// IsAuthorized - function to check the authorization and handle the function
func IsAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Authorization"] != nil {

			tokenString := strings.Split(r.Header["Authorization"][0], " ")[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				claims := token.Claims.(jwt.MapClaims)
				c := claims["userID"]
				cstr := fmt.Sprintf("%v", c)
				r.Header.Set("user_id", cstr)
				return []byte(cstr), nil
			})

			if err != nil {
				log.Println(err)
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

//validEmailID - function o validate email-id
func validEmailID(emailID string) bool {
	r := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return r.MatchString(emailID)
}

//validPhoneNo - function to validate phone number
func validPhoneNo(phoneNo string) bool {
	r := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return r.MatchString(phoneNo)
}

func renderHTML(w http.ResponseWriter, templateName string) error {
	dir := "views"

	templateData := GetLoginPageTemplateData{
		StaticDir: "/",
	}
	templates := make(map[string]*template.Template)

	if _, ok := templates[templateName]; !ok {
		tmpl, err := template.New(fmt.Sprintf("%s.tmpl.html", templateName)).ParseFiles(fmt.Sprintf("%s/%s.tmpl.html", dir, templateName))
		if err != nil {
		}
		templates[templateName] = tmpl
	}

	tmpl := templates[templateName]

	if err := tmpl.Execute(w, encode(templateData)); err != nil {

	}
	return nil
}

func encode(v interface{}) interface{} {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	var v2 interface{}
	if err := json.Unmarshal(data, &v2); err != nil {
		panic(err)
	}
	return v2
}
