package providers

type View struct {
	Data         []byte
	RefreshAfter int
}

type Provider interface {
	GetView(cursor string) ViewResponse
}

type ViewResponse struct {
	View       View
	Cursor     string
	NextCursor string
}
