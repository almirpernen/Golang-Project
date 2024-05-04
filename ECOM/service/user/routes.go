package user

import (
	"net/http"

	"github.com/almirpernen/types"
	"github.com/almirpernen/utils"
	"github.com/gorilla/mux"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/godoc/util"
	"gorm.io/gorm/utils"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter,r *http.Request) {

    // do something
}

func (h *Handler) handleRegister(w http.ResponseWriter,r *http.Request) {

    // do something

	//get JSON payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, payload) ; err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	//check if the user exists
	_, err. := h.store.GetUserByEmail(payload.Email)
	if err == nil
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email
		&s already exists", payload.Email))
		return
	}

	hashedPassword, err:= auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// if it doesnt we create the new user
	err = h.store.CreateUser(types.User{
		FirstName: 	payload.FirstName,
		LastName: 	payload.LastName,
		Email: 		payload.Email,
		Password: 	hashedPassword,
	})
	if err != nil {
		utils.WriteError(w,http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}