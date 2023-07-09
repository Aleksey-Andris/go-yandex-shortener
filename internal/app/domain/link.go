package domain

type Link struct {
	ident   string
	fulLink string
}

func (l *Link) Ident() string {
	return l.ident
}

func (l *Link) FulLink() string {
	return l.fulLink
}

func (l *Link) SetIdent(shortLink string) {
	l.ident = shortLink
}

func (l *Link) SetFulLink(fulLink string) {
	l.fulLink = fulLink
}
