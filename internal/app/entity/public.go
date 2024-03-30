package entity

type PageableRequest struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (p *PageableRequest) Init() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = 15
	}
}

type PageableResponse struct {
	Total int `json:"total"`
	List  any `json:"data"`
}
