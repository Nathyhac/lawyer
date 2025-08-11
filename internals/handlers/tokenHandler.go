package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Nathac/go-api/internals/store"
	"github.com/Nathac/go-api/internals/store/tokens"
	"github.com/Nathac/go-api/internals/utils"
)

type TokenHandler struct {
	UserStore  store.UserStore
	TokenStore store.TokenStore
	Logger     *log.Logger
}

type createRequestToken struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(userStore store.UserStore, tokenStore store.TokenStore) *TokenHandler {
	logger := log.New(os.Stdout, "INFO", log.Ldate|log.Ltime)
	return &TokenHandler{
		UserStore:  userStore,
		TokenStore: tokenStore,
		Logger:     logger,
	}
	
}

func (t *TokenHandler) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req createRequestToken
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		t.Logger.Printf("Error decoding the request:%v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal error"})
	}

	user, err := t.UserStore.GetUserUsername(req.Username)

	if err != nil {
		t.Logger.Printf("Error getting the user by username:%v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal error"})
	}

	passwordMatches, err := user.Hash_Password.Matches(req.Password)
	if err != nil {
		t.Logger.Printf("ERROR: passwordhash.matches. %v\n", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"Invalid": "internal server error"})
	}

	if !passwordMatches {
		t.Logger.Printf("ERROR: %v", err)
		utils.WriteJson(w, http.StatusBadRequest, utils.Envelop{"error": "invalid credentials"})
	}

	token, err := t.TokenStore.CreateToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		t.Logger.Printf("error creating a token %v\n", err)
		utils.WriteJson(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	utils.WriteJson(w, http.StatusOK, utils.Envelop{"auth-token": token})

}
