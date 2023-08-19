package domain

type Link struct {
	ID      int32  `json:"uuid" db:"id"`
	Ident   string `json:"short_url" db:"short_url"`
	FulLink string `json:"original_url" db:"original_url"`
}
