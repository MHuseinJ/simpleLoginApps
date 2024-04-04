package repository

import (
	"context"
)

func (r *Repository) GetAccountByPhone(ctx context.Context, phone string) (output Account, err error) {
	return
}
func (r *Repository) UpdateAccount(ctx context.Context, account Account) (output Account, err error) {
	return
}
func (r *Repository) CreateAccount(ctx context.Context, account Account) (output Account, err error) {
	return
}
func (r *Repository) CreateAuth(ctx context.Context, auth Auth) (output Auth, err error) { return }
func (r *Repository) GetAuthByAccountId(ctx context.Context, id int) (output Account, err error) {
	return
}
