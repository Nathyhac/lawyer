package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/Nathac/go-api/internals/store"
	"github.com/Nathac/go-api/internals/utils"
	"github.com/go-chi/chi/v5"
)

type registerUserRequest struct {
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Email        string `json:"email"`
	UserName     string `json:"username"`
	Phone_number string `json:"phone_number"`
	Password     string `json:"password"`
}

type UserHandler struct {
	UserStore store.UserStore
	Logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore) *UserHandler {

	logger := log.New(os.Stdout, "INFO", log.Ldate|log.Ltime)
	handler := &UserHandler{
		UserStore: userStore,
		Logger:    logger,
	}
	return handler
}

func (u *UserHandler) validateRegisterRequest(reg *registerUserRequest) error {
	if reg.First_name == "" {
		return errors.New("first name cannot be empty")
	}
	if reg.Last_name == "" {
		return errors.New("last name cannot be empty")
	}
	if reg.Email == "" {
		return errors.New("email cannot be empty")
	}

	emailregex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailregex.MatchString(reg.Email) {
		return errors.New("email format not correct")

	}
	if reg.Phone_number == "" {
		return errors.New("phone number cannot be empty")
	}
	if reg.Password == "" {
		return errors.New(" password cannot be empty")
	}

	return nil
}

func (uh *UserHandler) HandlerUserRegister(w http.ResponseWriter, r *http.Request) {
	var reg registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		uh.Logger.Printf("error in the request payload: %v", err)
		utils.WriteJson(w, http.StatusOK, utils.Envelop{"error": "invalid request payload"})
	}

	err = uh.validateRegisterRequest(&reg)
	if err != nil {
		uh.Logger.Printf("error in validating the payload: %v", err)
		utils.WriteJson(w, http.StatusOK, utils.Envelop{"error": "payload validation failed"})
		return
	}

	user := &store.User{
		First_name:   reg.First_name,
		Last_name:    reg.Last_name,
		UserName:     reg.UserName,
		Email:        reg.Email,
		Phone_number: reg.Phone_number,
	}

	err = user.Hash_Password.Set(reg.Password)
	if err != nil {
		uh.Logger.Printf("error hashing the password: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "there's some error with the server"})
		return
	}

	createdUser, err := uh.UserStore.CreateUser(user)
	if err != nil {
		uh.Logger.Printf("error in Creating the user: %v", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "cannot create the user"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelop{"User": createdUser})

}

func (u *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {

	username := chi.URLParam(r, "username")
	result, err := u.UserStore.GetUserUsername(username)
	if err != nil {
		u.Logger.Printf("error getting User using username:%v\n", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": "internal server error"})
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelop{"result": result})

}
