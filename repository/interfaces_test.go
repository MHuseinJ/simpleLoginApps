package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"reflect"
	"testing"
	"time"
)

var acccount = &Account{
	Id:       1,
	FullName: "Momo",
	Password: "Passw0rd!Hash",
	Phone:    "+628123456789",
}

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

func TestNewRepository(t *testing.T) {
	type args struct {
		opts NewRepositoryOptions
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepository(tt.args.opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_CreateAccount(t *testing.T) {
	db, dbMock := NewMock()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		account Account
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput Account
		wantErr    bool
		mock       func()
	}{
		{
			name: "create account",
			fields: fields{
				Db: db,
			},
			args: args{
				Account{
					FullName: acccount.FullName,
					Password: acccount.Password,
					Phone:    acccount.Phone,
				},
			},
			wantOutput: *acccount,
			wantErr:    false,
			mock: func() {
				dbMock.ExpectExec("INSERT INTO account").WithArgs(acccount.Phone, acccount.Password, acccount.FullName).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &Repository{
				Db: tt.fields.Db,
			}

			gotOutput, err := r.CreateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("CreateAccount() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_GetAccountByPhoneAndPassword(t *testing.T) {
	db, dbMock := NewMock()
	ctxMock, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx      *context.Context
		phone    string
		password string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput Account
		wantErr    bool
		mock       func()
	}{
		{
			name: "testing get account by phone and password",
			fields: fields{
				Db: db,
			},
			args: args{
				ctx:      &ctxMock,
				phone:    acccount.Phone,
				password: acccount.Password,
			},
			wantOutput: Account{
				Id:       acccount.Id,
				FullName: acccount.FullName,
				Phone:    acccount.Phone,
			},
			wantErr: false,
			mock: func() {
				queryStr := "SELECT\\s+fullname,\\s+id,\\s+phone\\s+FROM\\s+account\\s+WHERE\\s+phone\\s+=\\s+\\$1\\s+AND\\s+password\\s+=\\s+\\$2\n"

				rows := sqlmock.NewRows([]string{"fullname", "id", "phone"}).
					AddRow(acccount.FullName, acccount.Id, acccount.Phone)
				dbMock.ExpectQuery(queryStr).WithArgs(acccount.Phone, acccount.Password).WillReturnRows(rows)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			r := &Repository{
				Db: tt.fields.Db,
			}
			gotOutput, err := r.GetAccountByPhoneAndPassword(*tt.args.ctx, tt.args.phone, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByPhoneAndPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("GetAccountByPhoneAndPassword() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_GetAccountByToken(t *testing.T) {
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput Account
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			gotOutput, err := r.GetAccountByToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("GetAccountByToken() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_UpdateAccount(t *testing.T) {
	//db, dbMock := NewMock()
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		account Account
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput Account
		wantErr    bool
		mock       func()
	}{
		//{
		//	name: "testing update account by phone and password",
		//	fields: fields{
		//		Db: db,
		//	},
		//	args: args{
		//		account: *acccount,
		//	},
		//	wantOutput: Account{
		//		Id:       acccount.Id,
		//		FullName: acccount.FullName,
		//		Phone:    acccount.Phone,
		//	},
		//	wantErr: false,
		//	mock: func() {
		//		queryStr := "UPDATE account SET phone"
		//		dbMock.ExpectExec(queryStr).WithArgs(acccount.Phone, acccount.FullName, acccount.Id).WillReturnResult(sqlmock.NewResult(1, 1))
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			gotOutput, err := r.UpdateAccount(tt.args.account)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("UpdateAccount() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func TestRepository_UpdateLoginData(t *testing.T) {
	type fields struct {
		Db *sql.DB
	}
	type args struct {
		account Account
		jwt     string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOutput Account
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Db: tt.fields.Db,
			}
			gotOutput, err := r.UpdateLoginData(tt.args.account, tt.args.jwt)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateLoginData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("UpdateLoginData() gotOutput = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
