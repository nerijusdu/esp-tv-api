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

type PosthogSite struct {
	Title     string `json:"title"`
	ProjectId string `json:"projectId"`
	InsightId string `json:"insightId"`
}

type PosthogInsightResponse struct {
	Result []struct {
		Data   []int    `json:"data"`
		Labels []string `json:"labels"`
		Days   []string `json:"days"`
	} `json:"result"`
}

type BskyFeedResponse struct {
	Feed []struct {
		Post BskyPost `json:"post"`
	} `json:"feed"`
}

type BskyPost struct {
	Uri    string `json:"uri"`
	Cid    string `json:"cid"`
	Author struct {
		Did         string `json:"did"`
		Handle      string `json:"handle"`
		DisplayName string `json:"displayName"`
		Avatar      string `json:"avatar"`
	} `json:"author"`
	Record struct {
		Type      string `json:"$type"`
		CreatedAt string `json:"createdAt"`
		Embed     struct {
			Type   string `json:"$type"`
			Images []struct {
				Alt   string `json:"alt"`
				Image struct {
					Type     string `json:"$type"`
					MimeType string `json:"mimeType"`
					Size     int    `json:"size"`
					Ref      struct {
						Link string `json:"$link"`
					} `json:"ref"`
				} `json:"image"`
			} `json:"images"`
		} `json:"embed"`
		Text string `json:"text"`
	} `json:"record"`
	Embed struct {
		Type   string `json:"$type"`
		Images []struct {
			Thumb    string `json:"thumb"`
			Fullsize string `json:"fullsize"`
			Alt      string `json:"alt"`
		} `json:"images"`
	} `json:"embed"`
	ReplyCount  int `json:"replyCount"`
	LikeCount   int `json:"likeCount"`
	RepostCount int `json:"repostCount"`
	QuoteCount  int `json:"quoteCount"`
}

type WeatherResponse struct {
	Weather []struct {
		Id          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Sys struct {
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Name string `json:"name"`
}
