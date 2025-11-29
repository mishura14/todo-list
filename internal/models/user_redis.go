package models

type UserRedis struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}
