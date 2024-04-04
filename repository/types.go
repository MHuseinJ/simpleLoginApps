// This file contains types that are used in the repository layer.
package repository

type GetTestByIdInput struct {
	Id string
}

type GetTestByIdOutput struct {
	Name string
}

type Account struct {
	Id       int
	FullName string
	Phone    string
	Password string
}

type Login struct {
	Id           int
	AccountId    string
	Token        string
	successLogin int
}
