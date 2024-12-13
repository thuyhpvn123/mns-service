package request

type GetNamesRequest struct {
	Owner string `form:"owner" validate:"required"` 
}

type GetOwnerRequest struct {
	Name string `form:"name" validate:"required"` 
}

type VerifyDomainRequest struct {
	Domain string `json:"domain"` 
	Owner string `json:"owner"` 
	Label string `json:"label"` 
}