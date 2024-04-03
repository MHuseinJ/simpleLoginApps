package handler

import (
	"encoding/json"
	"fmt"
	"github.com/SawitProRecruitment/UserService/repository"
	"net/http"
	"strconv"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	var resp generated.BasicResponse
	idString := strconv.Itoa(params.Id)
	fmt.Println(idString)
	outPut, err := s.Repository.GetTestById(ctx.Request().Context(), repository.GetTestByIdInput{Id: idString})
	if err != nil {
		return err
	}
	resp.Message = fmt.Sprintf("Hello User %s", outPut.Name)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	//TODO implement me
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
