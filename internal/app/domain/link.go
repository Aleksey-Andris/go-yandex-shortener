// The domain package contains a business entityes.
package domain

// Link entity.
type Link struct {
	ID          int32  `json:"uuid" db:"id"`
	Ident       string `json:"short_url" db:"short_url"`
	FulLink     string `json:"original_url" db:"original_url"`
	UserID      int32  `json:"user_id" db:"user_id"`
	DeletedFlag bool   `json:"is_deleted" db:"is_deleted"`
}
