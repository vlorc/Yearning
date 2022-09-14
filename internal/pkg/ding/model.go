package ding

const (
	TYPE_MARKDOWN    = "markdown"
	TYPE_TEXT        = "text"
	TYPE_LINK        = "link"
	TYPE_ACTION_CARD = "actionCard"
	TYPE_FEED_CARD   = "FeedCard"
)

/**
docs: https://open.dingtalk.com/document/group/custom-robot-access
*/
type message struct {
	Content string `json:"content,omitempty"`

	Title    string `json:"title,omitempty"`
	Text     string `json:"text,omitempty"`
	Url      string `json:"messageUrl,omitempty"`
	ImageUrl string `json:"picUrl,omitempty"`

	SingleTitle       string   `json:"singleTitle,omitempty"`
	SingleURL         string   `json:"singleURL,omitempty"`
	ButtonOrientation string   `json:"btnOrientation,omitempty"` // 0：按钮竖直排列 1：按钮横向排列
	Button            []Button `json:"btns,omitempty"`

	Links []Link `json:"links,omitempty"`
}

type Button struct {
	Title string `json:"title"`
	Url   string `json:"actionURL"`
}

type Link struct {
	Title    string `json:"title,omitempty"`
	Url      string `json:"messageURL,omitempty"`
	ImageUrl string `json:"picURL,omitempty"`
}

type At struct {
	Mobiles []string `json:"atMobiles,omitempty"`
	UserIds []string `json:"atUserIds,omitempty"`
	All     bool     `json:"isAtAll"`
}

type result struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}
