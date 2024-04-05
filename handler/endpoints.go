package handler

import (
	"encoding/json"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) Login(ctx echo.Context) error {
	var resp generated.LoginResponse
	context := ctx.Request().Context()
	json_map := make(map[string]interface{})
	err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, generated.BasicResponse{Message: err.Error()})
	} else {
		phone := json_map["phone"].(string)
		password := json_map["password"].(string)
		var foundAccount repository.Account
		hashedPassword := CreateHash(password)
		fmt.Println(hashedPassword)
		foundAccount, err = s.Repository.GetAccountByPhoneAndPassword(context, phone, hashedPassword)
		if err != nil {
			fmt.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, generated.BasicResponse{Message: "account not found, reason: " + err.Error()})
		}
		token, err := createToken(foundAccount.FullName)
		if err != nil {
			fmt.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, generated.BasicResponse{Message: "create token error, reason: " + err.Error()})
		}
		foundAccount, err = s.Repository.UpdateLoginData(foundAccount, token)
		if err != nil {
			fmt.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, generated.BasicResponse{Message: "update login failed, reason: " + err.Error()})
		}
		resp.Message = fmt.Sprintf("success login for +%s", phone)
		resp.Token = token
		resp.Id = foundAccount.Id
		return ctx.JSON(http.StatusOK, resp)
	}
}

func (s *Server) GetProfile(ctx echo.Context, params generated.GetProfileParams) error {
	err := verifyToken(params.Authorization)
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusForbidden, generated.BasicResponse{Message: "Forbidden Access: " + err.Error()})
	}
	account, err := s.Repository.GetAccountByToken(ctx.Request().Context(), params.Authorization)

	return ctx.JSON(http.StatusOK, generated.ProfileResponse{
		Message:  "Success Retrieve Profile",
		Fullname: account.FullName,
		Phone:    account.Phone,
	})
}

func (s *Server) UpdateProfile(ctx echo.Context, params generated.UpdateProfileParams) error {
	err := verifyToken(params.Authorization)
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusForbidden, generated.BasicResponse{Message: "Forbidden Access: " + err.Error()})
	}
	account, err := s.Repository.GetAccountByToken(ctx.Request().Context(), params.Authorization)
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusNotFound, generated.BasicResponse{Message: err.Error()})
	}
	var resp generated.ProfileResponse
	json_map := make(map[string]interface{})
	err = json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, generated.BasicResponse{Message: err.Error()})
	} else {
		var phone string
		var fullname string
		if json_map["phone"] != nil {
			phone = json_map["phone"].(string)
		}
		if json_map["fullname"] != nil {
			fullname = json_map["fullname"].(string)
		}
		updatedAccount := repository.Account{Id: account.Id, FullName: fullname, Phone: phone}
		updatedAccount, err = s.Repository.UpdateAccount(updatedAccount)
		if err != nil {
			fmt.Println(err.Error())
			return ctx.JSON(http.StatusConflict, generated.BasicResponse{Message: err.Error()})
		}
		fmt.Println(fullname, phone)
		resp.Message = fmt.Sprintf("success updated profile id: %d", updatedAccount.Id)
		resp.Fullname = fullname
		resp.Phone = phone
		return ctx.JSON(http.StatusOK, resp)
	}
}

func (s *Server) Register(ctx echo.Context) error {
	var req generated.RegisterRequest
	var resp generated.RegisterResponse
	json_map := make(map[string]interface{})
	err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, generated.BasicResponse{Message: err.Error()})
	} else {
		req.Phone = json_map["phone"].(string)
		req.Password = json_map["password"].(string)
		req.Fullname = json_map["fullname"].(string)
		errors := validateRegisterRequest(req)
		if len(errors) > 0 {
			return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
				Message: "failed create account, bad request",
				Errors:  errors,
			})
		}
		hashedPassword := CreateHash(req.Password)
		fmt.Println(hashedPassword)
		account := repository.Account{
			Phone:    req.Phone,
			Password: hashedPassword,
			FullName: req.Fullname,
		}
		createdAccount, err := s.Repository.CreateAccount(account)
		if err != nil {
			fmt.Println(err.Error())
			return ctx.JSON(http.StatusBadRequest, generated.BasicResponse{Message: "failed create account, reason: " + err.Error()})
		}
		resp.Message = fmt.Sprintf("success create account for %s", account.Phone)
		resp.Id = createdAccount.Id
		return ctx.JSON(http.StatusOK, resp)
	}
}
