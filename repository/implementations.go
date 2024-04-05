package repository

import (
	"context"
	"fmt"
)

func (r *Repository) GetAccountByPhoneAndPassword(ctx context.Context, phone string, password string) (output Account, err error) {
	queryStr := `SELECT fullname, id, phone FROM account WHERE phone = $1 AND password = $2`
	err = r.Db.QueryRowContext(ctx, queryStr, phone, password).Scan(&output.FullName, &output.Id, &output.Phone)
	if err != nil {
		return Account{}, err
	}
	return output, err
}
func (r *Repository) UpdateAccount(account Account) (output Account, err error) {
	var queryStr string
	err = nil
	if account.Phone != "" && account.FullName != "" {
		queryStr = `UPDATE account SET phone = $1 ,fullname = $2 WHERE id = $3`
		_, err = r.Db.Exec(queryStr, account.Phone, account.FullName, account.Id)
	} else if account.Phone == "" {
		queryStr = `UPDATE account SET fullname = $1 WHERE id = $2`
		_, err = r.Db.Exec(queryStr, account.FullName, account.Id)
	} else if account.FullName == "" {
		queryStr = `UPDATE account SET phone = $1 WHERE id = $2`
		_, err = r.Db.Exec(queryStr, account.Phone, account.Id)
	}
	fmt.Println(queryStr)
	if err != nil {
		return account, err
	}
	return account, nil
}
func (r *Repository) CreateAccount(account Account) (output Account, err error) {
	queryStr := `INSERT INTO account (phone, password, fullname) 
VALUES ($1, $2, $3)`
	result, err := r.Db.Exec(queryStr, account.Phone, account.Password, account.FullName)
	if err != nil {
		return account, err
	}
	lastInsertedId, _ := result.LastInsertId()
	account.Id = int(lastInsertedId)
	return account, nil
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
func (r *Repository) GetAccountByToken(ctx context.Context, token string) (output Account, err error) {
	queryStr := `SELECT a.id, a.phone, a.fullname FROM account a JOIN login l ON a.id = l.account_id WHERE l.token = $1`
	err = r.Db.QueryRowContext(ctx, queryStr, token).Scan(&output.Id, &output.Phone, &output.FullName)
	if err != nil {
		return Account{}, err
	}
	return output, err
}
