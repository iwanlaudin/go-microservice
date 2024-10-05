package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iwanlaudin/go-microservice/pkg/common/api"
	"github.com/iwanlaudin/go-microservice/pkg/common/helpers"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/dto/request"
	"github.com/iwanlaudin/go-microservice/services/authentication/internal/service"
)

type AuthHandler struct {
	AuthService service.AuthService
	Validate    *validator.Validate
}

func NewAuthHandler(authService service.AuthService, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Validate:    validate,
	}
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	request := request.CreateUserRequest{}

	helpers.ReadFromRequestBody(r, &request)
	if err := h.Validate.Struct(request); err != nil {
		api.NewAppResponse("Invalid parameter", http.StatusBadRequest).ValidationErr(w, err)
		return
	}

	response, err := h.AuthService.Create(r.Context(), request)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(w)
		return
	}

	api.NewAppResponse("Registration successfully", http.StatusOK).Ok(w, response)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	request := request.SignInRequest{}

	helpers.ReadFromRequestBody(r, &request)
	if err := h.Validate.Struct(request); err != nil {
		api.NewAppResponse("Invalid parameter", http.StatusBadRequest).ValidationErr(w, err)
		return
	}

	user, err := h.AuthService.SignIn(r.Context(), &request)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).ValidationErr(w, err)
		return
	}

	api.NewAppResponse("Login successfully", http.StatusOK).Ok(w, user)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	request := request.RefreshTokenRequest{}

	helpers.ReadFromRequestBody(r, request)
	if err := h.Validate.Struct(request); err != nil {
		api.NewAppResponse("Invalid parameter", http.StatusBadRequest).ValidationErr(w, err)
		return
	}

	userToken, err := h.AuthService.RefreshToken(r.Context(), &request)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).ValidationErr(w, err)
		return
	}
	api.NewAppResponse("Refresh token successfully", http.StatusOK).Ok(w, userToken)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	idStr := api.UserIDFromContext(r.Context())

	userId, err := helpers.ConvertUserIDToUUID(idStr)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusInternalServerError).Err(w)
		return
	}

	userResponse, err := h.AuthService.FindUserById(r.Context(), userId)
	if err != nil {
		api.NewAppResponse(err.Error(), http.StatusBadRequest).Err(w)
		return
	}
	api.NewAppResponse("Successfully get the user", http.StatusOK).Ok(w, userResponse)
}
