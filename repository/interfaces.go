// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	GetAccountByPhoneAndPassword(ctx context.Context, phone string, password string) (output Account, err error)
	UpdateAccount(ctx context.Context, account Account) (output Account, err error)
	CreateAccount(ctx context.Context, account Account) (output Account, err error)
	UpdateLoginData(ctx context.Context, account Account, jwt string) (output Account, err error)
	GetAuthByAccountId(ctx context.Context, id int) (output Account, err error)
}
