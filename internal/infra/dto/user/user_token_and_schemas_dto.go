package userdto

type UserTokenDTO struct {
	User        UserDTO `json:"user"`
	AccessToken string  `json:"access_token"`
}
