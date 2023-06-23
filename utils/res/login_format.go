package res

import "point-of-sale/app/model"

type SetLoginFormat struct {
	Username  string `json:"username"`
	UsersCode string `json:"user_code"`
	Token     string `json:"token"`
}

func TransformLoginResponse(request model.User, token string) SetLoginFormat {
	return SetLoginFormat{
		Username:  request.Username,
		UsersCode: request.UserCode,
		Token:     token,

	}
}


