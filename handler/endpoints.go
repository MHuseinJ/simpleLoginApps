package handler

import (
	"encoding/json"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) Login(ctx echo.Context) error {
	var resp generated.BasicResponse
	json_map := make(map[string]interface{})
	err := json.NewDecoder(ctx.Request().Body).Decode(&json_map)
	if err != nil {
		return err
	} else {
		phone := json_map["phone"]
		password := json_map["password"]
		resp.Message = fmt.Sprintf("Login %d Password %s", phone, password)
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
