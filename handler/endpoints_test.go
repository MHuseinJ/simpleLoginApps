package handler

import (
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
				mockRepo.EXPECT().GetAccountByPhoneAndPassword(ctxMock.Request().Context(), "+62812345678", createHash("password")).Return(repository.Account{}, nil)
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
					Password: createHash("Password!1"),
				}).Return(repository.Account{}, nil)
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
					`{"phone":"+62812345678","password":"password","fullname":"john due"}`))
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
	type fields struct {
		Repository repository.RepositoryInterface
	}
	type args struct {
		ctx    echo.Context
		params generated.UpdateProfileParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Repository: tt.fields.Repository,
			}
			if err := s.UpdateProfile(tt.args.ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_createHash(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createHash(tt.args.password); got != tt.want {
				t.Errorf("createHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createToken(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createToken(tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("createToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("createToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateRegisterRequest(t *testing.T) {
	type args struct {
		request generated.RegisterRequest
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateRegisterRequest(tt.args.request); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateRegisterRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_verifyToken(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyToken(tt.args.tokenString); (err != nil) != tt.wantErr {
				t.Errorf("verifyToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
