// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type Account struct {
	Id       string
	FullName string
	Phone    string
	Password string
}

type Auth struct {
	Id        string
	AccountId string
	Token     string
	Expired   string
}
