package messagex

type Type string

const (
	TYPE_MARKDOWN    Type = "markdown"
	TYPE_TEXT        Type = "text"
	TYPE_HTML        Type = "html"
	TYPE_IMAGE       Type = "image"
	TYPE_FILE        Type = "file"
	TYPE_LINK        Type = "link"
	TYPE_ACTION_CARD Type = "actionCard"
	TYPE_FEED_CARD   Type = "FeedCard"
)

type Message struct {
	Type     Type     `json:"type"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	Link     string   `json:"link"`
	ImageUrl string   `json:"imageUrl"`
	Target   Target   `json:"target"`
	Files    []File   `json:"files"`
	Params   []Param  `json:"params"`
	App      *AppInfo `json:"app"`
	Sources  []Source `json:"sources"`
}

type AppInfo struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	IconUrl string `json:"iconUrl"`
}

type Source struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	OpenId string `json:"openId"`
	Mobile string `json:"mobile"`
	Kind   string `json:"kind"`
	Title  string `json:"title"`
	Phrase string `json:"phrase"`
}

type Target struct {
	All     bool     `json:"all"`
	Emails  []string `json:"emails"`
	OpenIds []string `json:"openIds"`
	Mobiles []string `json:"mobiles"`
}

type Param struct {
	Name  string
	Value string
	Color string
}

type File struct {
	Name   string
	Data   []byte
	Inline bool
}
