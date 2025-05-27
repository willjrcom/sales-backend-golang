package companydto

type UserTokenDTO struct {
	User    UserDTO `json:"user"`
	IDToken string  `json:"id_token"`
}
