package handler

import (
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	type args struct {
		opts NewServerOptions
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{
			name: "testing constructor for handler",
			args: args{
				opts: NewServerOptions{Repository: mockRepo},
			},
			want: &Server{
				Repository: mockRepo,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServer(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctxMock := e.NewContext(req, rec)
	token, _ := createToken("token")

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx    echo.Context
		params generated.GetProfileParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		mock    func()
		wantErr bool
	}{
		{
			name: "testing get profile",
			args: args{
				ctx: ctxMock,
				params: generated.GetProfileParams{
					Authorization: token,
				},
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				mockRepo.EXPECT().GetAccountByToken(ctxMock.Request().Context(), token).Return(repository.Account{}, nil)
			},
		},
		{
			name: "testing get profile with false token",
			args: args{
				ctx: ctxMock,
				params: generated.GetProfileParams{
					Authorization: token + "false",
				},
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				Repository: tt.fields.Repository,
			}
			if err := s.GetProfile(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("GetProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var ctxMock echo.Context
	//var token string
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()
	e := echo.New()

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx *echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "testing login with correct credential",
			args: args{
				ctx: &ctxMock,
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(
					`{"phone":"+62812345678","password":"password"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
				mockRepo.EXPECT().GetAccountByPhoneAndPassword(ctxMock.Request().Context(), "+62812345678", CreateHash("password")).Return(repository.Account{}, nil)
				mockRepo.EXPECT().UpdateLoginData(repository.Account{}, gomock.Any()).Return(repository.Account{}, nil)
			},
		},
		{
			name: "testing login with uncorrected json format",
			args: args{
				ctx: &ctxMock,
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPost, "/profile", strings.NewReader(
					`{"phone":"+62812345678","password":"password",}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				Repository: tt.fields.Repository,
			}
			if err := s.Login(*tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var ctxMock echo.Context
	//var token string
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()
	e := echo.New()

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx *echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "register with correct credential",
			args: args{
				ctx: &ctxMock,
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(
					`{"phone":"+62812345678","password":"Password!1","fullname":"john due"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
				mockRepo.EXPECT().CreateAccount(repository.Account{
					FullName: "john due",
					Phone:    "+62812345678",
					Password: CreateHash("Password!1"),
				}).Return(repository.Account{}, nil)
			},
		},
		{
			name: "register with existing phone number",
			args: args{
				ctx: &ctxMock,
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(
					`{"phone":"+62812345678","password":"Password!1","fullname":"john due"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
				mockRepo.EXPECT().CreateAccount(repository.Account{
					FullName: "john due",
					Phone:    "+62812345678",
					Password: CreateHash("Password!1"),
				}).Return(repository.Account{}, fmt.Errorf("duplicate constraint for phone"))
			},
		},
		{
			name: "register with unqualified password",
			args: args{
				ctx: &ctxMock,
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(
					`{"phone":"+62812345678","password":"pass","fullname":"john due"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
			},
		},
		{
			name: "register with unqualified phone and fullname",
			args: args{
				ctx: &ctxMock,
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(
					`{"phone":"0812345678","password":"password","fullname":"jo"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				Repository: tt.fields.Repository,
			}
			if err := s.Register(*tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServer_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	var rec *httptest.ResponseRecorder
	var req *http.Request
	var ctxMock echo.Context
	mockRepo := repository.NewMockRepositoryInterface(ctrl)
	defer ctrl.Finish()
	e := echo.New()
	token, _ := createToken("token")

	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx    *echo.Context
		params generated.UpdateProfileParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mock    func()
	}{
		{
			name: "update profile with phone and fullname",
			args: args{
				ctx: &ctxMock,
				params: generated.UpdateProfileParams{
					Authorization: token,
				},
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(
					`{"phone":"+62812345678","fullname":"john due"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				req.Header.Set(echo.HeaderAuthorization, token)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
				mockRepo.EXPECT().GetAccountByToken(ctxMock.Request().Context(), token).Return(repository.Account{Id: 1}, nil)
				mockRepo.EXPECT().UpdateAccount(repository.Account{
					Id:       1,
					FullName: "john due",
					Phone:    "+62812345678",
				}).Return(repository.Account{}, nil)
			},
		},
		{
			name: "update profile with phone and duplicated",
			args: args{
				ctx: &ctxMock,
				params: generated.UpdateProfileParams{
					Authorization: token,
				},
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(
					`{"phone":"+62812345678"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				req.Header.Set(echo.HeaderAuthorization, token)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
				mockRepo.EXPECT().GetAccountByToken(ctxMock.Request().Context(), token).Return(repository.Account{Id: 1}, nil)
				mockRepo.EXPECT().UpdateAccount(repository.Account{
					Id:    1,
					Phone: "+62812345678",
				}).Return(repository.Account{}, fmt.Errorf("duplicate phone number"))
			},
		},
		{
			name: "update profile with fullname only",
			args: args{
				ctx: &ctxMock,
				params: generated.UpdateProfileParams{
					Authorization: token,
				},
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(
					`{"fullname":"john due updated"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				req.Header.Set(echo.HeaderAuthorization, token)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
				mockRepo.EXPECT().GetAccountByToken(ctxMock.Request().Context(), token).Return(repository.Account{Id: 1}, nil)
				mockRepo.EXPECT().UpdateAccount(repository.Account{
					Id:       1,
					FullName: "john due updated",
				}).Return(repository.Account{}, fmt.Errorf("duplicate phone number"))
			},
		},
		{
			name: "update profile with unexist profile",
			args: args{
				ctx: &ctxMock,
				params: generated.UpdateProfileParams{
					Authorization: token,
				},
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(
					`{"phone":"+62812345678"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				req.Header.Set(echo.HeaderAuthorization, token)
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
				mockRepo.EXPECT().GetAccountByToken(ctxMock.Request().Context(), token).Return(repository.Account{}, fmt.Errorf("account not found"))
			},
		},
		{
			name: "update profile with unvalid token",
			args: args{
				ctx: &ctxMock,
				params: generated.UpdateProfileParams{
					Authorization: "unvalidtoken",
				},
			},
			fields: fields{
				Repository: mockRepo,
			},
			wantErr: false,
			mock: func() {
				req = httptest.NewRequest(http.MethodPut, "/profile", strings.NewReader(
					`{"phone":"+62812345678"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				req.Header.Set(echo.HeaderAuthorization, "unvalidtoken")
				rec = httptest.NewRecorder()
				ctxMock = e.NewContext(req, rec)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				Repository: tt.fields.Repository,
			}
			if err := s.UpdateProfile(*tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
