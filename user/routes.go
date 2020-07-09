package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/unrolled/render"
)

type userHandler struct {
	db         *sqlx.DB
	repository UserRepository
	render     *render.Render
}

func NewUserHandler(db *sqlx.DB) *userHandler {
	return &userHandler{
		db:         db,
		repository: UserRepository{db: *db},
		render:     render.New(),
	}
}

func getUserFromRequest(w http.ResponseWriter, r *http.Request) (*User, error) {
	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	if err := user.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	return user, nil
}

func (u *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	if user, err := u.repository.Get(chi.URLParam(r, "id")); err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		u.render.JSON(w, http.StatusOK, user)
	}
}

func (u *userHandler) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(w, r)
	if err != nil {
		return
	}

	newUser, err := u.repository.CreateUser(*user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	u.render.JSON(w, http.StatusCreated, newUser)
}

func (u *userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(w, r)
	if err != nil {
		return
	}

	id := chi.URLParam(r, "id")
	if user.ID != "" && user.ID != id {
		http.Error(w, "User id is not equal URL ID", http.StatusBadRequest)
		return
	}

	user.ID = id
	if err = u.repository.Update(*user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u *userHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	if err := u.repository.Delete(chi.URLParam(r, "id")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func RegisterRoutes(db *sqlx.DB, router chi.Router) {
	handler := NewUserHandler(db)

	router.Route("/user", func(r chi.Router) {
		r.Get("/{id}", handler.getUser)
		r.Post("/", handler.createUser)
		r.Delete("/{id}", handler.deleteUser)
		r.Put("/{id}", handler.updateUser)
	})
}
