package repository

import (
	"context"
)

func (r *Repository) GetAccountByPhoneAndPassword(ctx context.Context, phone string, password string) (output Account, err error) {
	queryStr := `SELECT full_name, id, phone FROM account WHERE phone = $1 AND password = $2`
	err = r.Db.QueryRowContext(ctx, queryStr, phone, password).Scan(&output.FullName, &output.Id, &output.Phone)
	if err != nil {
		return Account{}, err
	}
	return output, err
}
func (r *Repository) UpdateAccount(ctx context.Context, account Account) (output Account, err error) {
	return
}
func (r *Repository) CreateAccount(ctx context.Context, account Account) (output Account, err error) {
	return
}
func (r *Repository) UpdateLoginData(account Account, jwt string) (output Account, err error) {
	queryStr := `INSERT INTO login (account_id, success_login, token) 
VALUES ($1, 1, $2) 
ON CONFLICT (account_id) DO UPDATE SET success_login = login.success_login + 1, token = $2`
	_, err = r.Db.Exec(queryStr, account.Id, jwt)
	if err != nil {
		return account, err
	}
	return account, nil
}
func (r *Repository) GetAuthByAccountId(ctx context.Context, id int) (output Account, err error) {
	return
}
