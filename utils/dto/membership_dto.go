package dto

type AddMembershipRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
	BirthDay string `json:"birth_day"`
}

type AddPointRequest struct {
	ID               int `json:"id"`
	TotalTransaction int `json:"total_transaction"`
}

type EditMembershipRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    int    `json:"phone"`
	BirthDay string `json:"birth_day"`
}

type SearchMembershipRequest struct {
	MemberCode string `json:"member_code"`
}