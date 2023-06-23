package dto

type RequestAuthInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestAuthEmployee struct {
	MemberID string `json:"member_id"`
	Password string `json:"password"`
}
