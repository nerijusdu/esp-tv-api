package providers

type View struct {
	Data         []byte
	RefreshAfter int
}

type Provider interface {
	GetView(cursor string) (ViewResponse, error)
	Init(config any) error
}

type ViewResponse struct {
	View       View
	Cursor     string
	NextCursor string
}
