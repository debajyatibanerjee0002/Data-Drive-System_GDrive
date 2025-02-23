package entities

import (
	"encoding/json"
	"io"
)

type SignUpResponseEntity struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type SignUpReqestEntity struct {
	UserId   int    `json:"userid,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserName string `json:"username,omitempty"`
}

type LoginRequestEntity struct {
	EmailOrUsername string `json:"loginid"`
	Password        string `json:"password"`
}

type LoginResponseEntity struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func (w *SignUpReqestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
func (w *LoginRequestEntity) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(w)
}
