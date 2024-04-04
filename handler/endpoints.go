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
		hashedPassword := createHash(password)
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
	//Todo Implement Get Profile
	return nil
}

func (s *Server) UpdateProfile(ctx echo.Context, params generated.UpdateProfileParams) error {
	//Todo Implement Update Profile
	return nil
}

func (s *Server) Register(ctx echo.Context) error {
	//Todo Implement Register Profile
	return nil
}
