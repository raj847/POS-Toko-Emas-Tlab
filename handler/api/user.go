package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"tim-b/entity"
	"tim-b/service"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type UserAPI struct {
	userService *service.UserService
}

func NewUserAPI(
	userService *service.UserService,
) *UserAPI {
	return &UserAPI{
		userService: userService,
	}
}

func (u *UserAPI) AddAnggota(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Username == "" || user.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("register data is empty"))
		return
	}

	usernameClaims := r.Context().Value("username").(string)
	user.CreatedBy = usernameClaims

	eUser, err := u.userService.AddAnggota(r.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrPasswordInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrEmailInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id": eUser.ID,
		"message": "add anggota success",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (u *UserAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println(err)
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Username == "" || user.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("register data is empty"))
		return
	}

	eUser, err := u.userService.AddAnggota(r.Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrPasswordInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrEmailInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id": eUser.ID,
		"message": "add anggota success",
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (u *UserAPI) UserLogin(w http.ResponseWriter, r *http.Request) {
	var userReq entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	if userReq.Username == "" || userReq.Password == "" {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("username or password is empty"))
		return
	}

	eUser, err := u.userService.LoginUser(r.Context(), userReq)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserPasswordDontMatch) {
			WriteJSON(w, http.StatusNotFound, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	expiresAt := time.Now().Add(5 * time.Hour)
	claims := entity.Claims{
		UserID:   eUser.ID,
		Username: eUser.Username,
		Role:     eUser.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, _ := token.SignedString([]byte("rahasia-perusahaan"))

	response := map[string]any{
		"data":        eUser,
		"tokenCookie": tokenString,
	}

	WriteJSON(w, http.StatusOK, response)
}

func (u *UserAPI) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query()
	adder, foundAdder := user["created_by"]
	userID, foundUserId := user["user_id"]
	nama := r.URL.Query().Get("nama")

	if foundAdder {
		creator := adder[0]
		anggota, err := u.userService.GetUsersbyAdder(r.Context(), creator)
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		WriteJSON(w, http.StatusOK, anggota)
		return
	}

	if foundUserId {
		userId, _ := strconv.Atoi(userID[0])
		anggota, err := u.userService.GetUsersbyID(r.Context(), uint(userId))
		if err != nil {
			WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
			return
		}

		WriteJSON(w, http.StatusOK, anggota)
		return
	}

	hasilSearch := []entity.User{}
	listUser, err := u.userService.GetAllUsers(r.Context())
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}
	if nama != "" {
		for _, item := range listUser {
			lowerName := strings.ToLower(item.NamaLengkap)
			lowerQuery := strings.ToLower(nama)
			if strings.Contains(lowerName, lowerQuery) {
				hasilSearch = append(hasilSearch, item)
			}
		}
		WriteJSON(w, http.StatusOK, hasilSearch)
		return
	}

	WriteJSON(w, http.StatusOK, listUser)
}

func (u *UserAPI) UpdateAnggota(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid decode json"))
		return
	}

	usernameClaims := r.Context().Value("username").(string)
	user.UpdatedBy = usernameClaims

	id := r.URL.Query().Get("user_id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse("invalid user id"))
		return
	}

	t, _ := time.Parse("2006-01-02", user.TanggalMasukStr)

	eUser, err := u.userService.UpdateUser(r.Context(), entity.User{
		Model: gorm.Model{
			ID: uint(idInt),
		},
		NamaLengkap:  user.NamaLengkap,
		Username:     user.Username,
		Password:     user.Password,
		NoHp:         user.NoHp,
		Email:        user.Email,
		TanggalMasuk: t,
		Status:       entity.Statusx(user.Status),
		CreatedBy:    user.CreatedBy,
		UpdatedBy:    user.UpdatedBy,
		Role:         user.Role,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			WriteJSON(w, http.StatusConflict, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrPasswordInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrUserInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		} else if errors.Is(err, service.ErrEmailInvalid) {
			WriteJSON(w, http.StatusBadRequest, entity.NewErrorResponse(err.Error()))
			return
		}

		WriteJSON(w, http.StatusInternalServerError, entity.NewErrorResponse("error internal server"))
		return
	}

	response := map[string]any{
		"user_id": eUser.ID,
		"message": "update anggota success",
	}

	WriteJSON(w, http.StatusCreated, response)
}
