package domain

type Link struct {
	shortLink string
	fulLink   string
}

func (l *Link) ShortLink() string {
	return l.shortLink
}

func (l *Link) FulLink() string {
	return l.fulLink
}

func (l *Link) SetShortLink(shortLink string) {
	l.shortLink = shortLink
}

func (l *Link) SetFulLink(fulLink string) {
	l.fulLink = fulLink
}
