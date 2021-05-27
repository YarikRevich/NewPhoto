package exceptions

import "errors"

const (
	LOGIN_ERROR       = "0"
	LOGOUT_ERROR      = "1"
	REGISTRAION_ERROR = "2"
)

//Returns error message gotten from passed error code
func GetErrorMessageByCode(c string) (string, error) {
	switch c {
	case LOGIN_ERROR:
	case LOGOUT_ERROR:
	case REGISTRAION_ERROR:
	}
	return "", errors.New("there is no such code")
}
