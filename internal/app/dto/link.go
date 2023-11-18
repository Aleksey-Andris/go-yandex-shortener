package dto

// LinkReq -  a structure describing the body of a request to shorten one link.
type LinkReq struct {
	// URL - original URL.
	URL string `json:"url"`
}

// LinkReq -  a structure describing the body of a response to shorten one link.
type LinkRes struct {
	// Result - shortened  URL including server address.
	Result string `json:"result"`
}

// LinkReq -  a structure describing the body of a request to shorten many links.
type LinkListReq struct {
	// CorrelationID - link ident within this http request.
	CorrelationID string `json:"correlation_id"`
	// OriginalURL - original URL.
	OriginalURL string `json:"original_url"`
}

// LinkReq -  a structure describing the body of a response to shorten many links.
type LinkListRes struct {
	// CorrelationID - link ident within this http request.
	CorrelationID string `json:"correlation_id"`
	// ShortURL - shortened  URL including server address.
	ShortURL string `json:"short_url"`
}

// LinkReq -  a structure describing the body of a response to user's links.
type LinkListByUserIDRes struct {
	// OriginalURL - original URL.
	OriginalURL string `json:"original_url" db:"original_url"`
	// ShortURL - shortened  URL including server address.
	ShortURL string `json:"short_url" db:"short_url"`
}
