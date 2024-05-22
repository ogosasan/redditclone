package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/session"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/repo/user"
	"gitlab.vk-golang.com/vk-golang/lectures/05_web_app/99_hw/redditclone/pkg/utils"

	"io"
	"net/http"

	"go.uber.org/zap"
	"html/template"
)

type UserHandler struct {
	Tmpl     *template.Template
	Logger   *zap.SugaredLogger
	UserRepo user.UserRepo
	Sessions *session.SessionsManager
}

type LoginRegisterRequestBody struct {
	Password string `json:"password" valid:"required,length(8|255)"`
	Username string `json:"username" valid:"required,matches(^[a-zA-Z0-9_]+$)"`
}

type loginRegisterResponseBody struct {
	Token string `json:"token"`
}

func (lr *LoginRegisterRequestBody) Validate() []string {
	_, err := govalidator.ValidateStruct(lr)
	errors := make([]string, 0)
	if err == nil {
		return errors
	}
	if errs, ok := err.(govalidator.Errors); ok {
		for _, e := range errs {
			errors = append(errors, e.Error())
		}
	}
	return errors
}

func (h *UserHandler) HandleGetToken(w http.ResponseWriter, user *user.User) {
	token, err := h.Sessions.CreateSession(user)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in session creation: %s"}`, err)
		utils.MakeResponse(h.Logger, w, []byte(errText), 500)
		return
	}
	response := loginRegisterResponseBody{token}
	responseJSON, err := json.Marshal(&response)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "error in coding response: %s"}`, err)
		utils.MakeResponse(h.Logger, w, []byte(errText), 500)
		return
	}
	h.Logger.Infof("Create new token %s", token)
	utils.MakeResponse(h.Logger, w, responseJSON, 200)
}

func RequestUserHandler(logger *zap.SugaredLogger, w http.ResponseWriter, r *http.Request) (*LoginRegisterRequestBody, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "Request can not be readed %s"}`, err))
		utils.MakeResponse(logger, w, textErr, 500)
		return nil, err
	}
	bodyUnmarshal := &LoginRegisterRequestBody{}
	err = json.Unmarshal(body, bodyUnmarshal)
	if err != nil {
		textErr := []byte(fmt.Sprintf(`{"message": "The request cannot be unmarshaled %s"}`, err))
		utils.MakeResponse(logger, w, textErr, 400)
		return nil, err
	}
	if validationErrors := bodyUnmarshal.Validate(); len(validationErrors) != 0 {
		errorsJSON, err := json.Marshal(validationErrors)
		if err != nil {
			errText := fmt.Sprintf(`{"message": "Validation errors can not be marshaled: %s"}`, err)
			utils.MakeResponse(logger, w, []byte(errText), 500)
			return nil, err
		}
		logger.Errorf("Validation is not success: %s", err)
		utils.MakeResponse(logger, w, errorsJSON, 500)
		return nil, err
	}
	return bodyUnmarshal, nil
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestData, err := RequestUserHandler(h.Logger, w, r)
	if err != nil || requestData == nil {
		errText := fmt.Sprintf(`{"message": "Error in receiving data: %s"}`, err)
		utils.MakeResponse(h.Logger, w, []byte(errText), 500)
		return
	}
	getLogged, err := h.UserRepo.Login(requestData.Username, requestData.Password)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "Error in register: %s"}`, err)
		utils.MakeResponse(h.Logger, w, []byte(errText), 500)
		return
	}
	h.HandleGetToken(w, getLogged)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestData, err := RequestUserHandler(h.Logger, w, r)
	if err != nil || requestData == nil {
		errText := fmt.Sprintf(`{"message": "Error in receiving data: %s"}`, err)
		utils.MakeResponse(h.Logger, w, []byte(errText), 500)
		return
	}
	newUser, err := h.UserRepo.Register(requestData.Username, requestData.Password)
	if err != nil {
		errText := fmt.Sprintf(`{"message": "Error in register: %s"}`, err)
		utils.MakeResponse(h.Logger, w, []byte(errText), 500)
		return
	}
	h.HandleGetToken(w, newUser)
}
