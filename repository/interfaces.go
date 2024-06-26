// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetAccountByPhoneAndPassword(ctx context.Context, phone string, password string) (output Account, err error)
	UpdateAccount(account Account) (output Account, err error)
	CreateAccount(account Account) (output Account, err error)
	UpdateLoginData(account Account, jwt string) (output Account, err error)
	GetAccountByToken(ctx context.Context, token string) (output Account, err error)
}
