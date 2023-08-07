package domain

type Link struct {
	ID      int32  `json:"uuid"`
	Ident   string `json:"short_url"`
	FulLink string `json:"original_url"`
}
