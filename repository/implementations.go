package repository

import (
	"context"
)

func (r *Repository) GetAccountByPhoneAndPassword(ctx context.Context, phone string, password string) (output Account, err error) {
	return
}
func (r *Repository) UpdateAccount(ctx context.Context, account Account) (output Account, err error) {
	return
}
func (r *Repository) CreateAccount(ctx context.Context, account Account) (output Account, err error) {
	return
}
func (r *Repository) UpdateLoginData(ctx context.Context, account Account, jwt string) (output Account, err error) {
	return
}
func (r *Repository) GetAuthByAccountId(ctx context.Context, id int) (output Account, err error) {
	return
}
