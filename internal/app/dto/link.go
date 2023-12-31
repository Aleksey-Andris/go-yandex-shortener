package dto

type LinkReq struct {
	URL string `json:"url"`
}

type LinkRes struct {
	Result string `json:"result"`
}

type LinkListReq struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type LinkListRes struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type LinkListByUserIDRes struct {
	OriginalURL string `json:"original_url" db:"original_url"`
	ShortURL    string `json:"short_url" db:"short_url"`
}
