package httpsrv

import (
	"encoding/json"
	"net/http"

	"github.com/dieklingel/core/internal/core"
	"github.com/dieklingel/core/internal/slice"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func buildUserRoutes(service *HttpService, router *mux.Router) {
	type User struct {
		Username string
	}

	router.Methods("GET").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users := slice.Map[core.User, User](service.UserService.Users(), func(user core.User) User {
			return User{
				Username: user.Username,
			}
		})

		payload, err := json.Marshal(users)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(payload)
	})

	router.Methods("PUT").Path("").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Request struct {
			Username string
			Password string
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if len(req.Username) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("the username can not be empty"))
			return
		}
		if len(req.Password) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("the password can not be empty"))
			return
		}
		if existingUser := service.UserService.GetUserByUsername(req.Username); existingUser != nil {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("a user with this username already exists"))
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}

		user := core.User{
			Username:     req.Username,
			PasswordHash: string(passwordHash),
		}
		service.UserService.SaveUser(&user)
		json.NewEncoder(w).Encode(User{
			Username: user.Username,
		})
	})

	router.Methods("GET").Path("/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["id"]

		user := service.UserService.GetUserByUsername(username)
		if user == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(User{
			Username: user.Username,
		})
	})

	router.Methods("DELETE").Path("/{id}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		username := vars["id"]

		user := service.UserService.GetUserByUsername(username)
		if user != nil {
			service.UserService.RemoveUser(user)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
