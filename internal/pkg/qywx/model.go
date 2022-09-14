package qywx

const (
	TYPE_MARKDOWN      = "markdown"
	TYPE_TEXT          = "text"
	TYPE_IMAGE         = "image"
	TYPE_NEWS          = "news"
	TYPE_FILE          = "file"
	TYPE_TEMPLATE_CARD = "template_card"
)

/**
docs: https://developer.work.weixin.qq.com/document/path/91770
*/
type message struct {
	Type         string               `json:"msgtype"`
	Text         *messageText         `json:"text,omitempty"`
	Markdown     *messageText         `json:"markdown,omitempty"`
	Image        *messageImage        `json:"image,omitempty"`
	News         *messageNews         `json:"news,omitempty"`
	File         *messageFile         `json:"file,omitempty"`
	TemplateCard *messageTemplateCard `json:"template_card,omitempty"`
}

type messageText struct {
	Content string `json:"content,omitempty"`

	UserIds []string `json:"mentioned_list,omitempty"`
	Mobiles []string `json:"mentioned_mobile_list,omitempty"`
}

type messageImage struct {
	Base64 string `json:"base64,omitempty"`
	Md5    string `json:"md5,omitempty"`
}

type messageNews struct {
	Articles []article `json:"articles,omitempty"`
}

type article struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
	ImageUrl    string `json:"picurl,omitempty"`
}

type messageFile struct {
	MediaId string `json:"media_id,omitempty"`
}

type messageTemplateCardSource struct {
	IconURL   string `json:"icon_url,omitempty"`
	Desc      string `json:"desc,omitempty"`
	DescColor int    `json:"desc_color,omitempty"`
}

type messageTemplateCardContent struct {
	Title string `json:"title,omitempty"`
	Desc  string `json:"desc,omitempty"`
}

type messageTemplateCardQuoteArea struct {
	Type      int    `json:"type,omitempty"`
	URL       string `json:"url,omitempty"`
	Appid     string `json:"appid,omitempty"`
	Pagepath  string `json:"pagepath,omitempty"`
	Title     string `json:"title,omitempty"`
	QuoteText string `json:"quote_text,omitempty"`
}

type messageTemplateCardCardAction struct {
	Type     int    `json:"type,omitempty"`
	URL      string `json:"url,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}

type messageTemplateCardHorizontalContent struct {
	Keyname string `json:"keyname,omitempty"`
	Value   string `json:"value,omitempty"`
	Type    int    `json:"type,omitempty"`
	URL     string `json:"url,omitempty"`
	MediaID string `json:"media_id,omitempty"`
}

type messageTemplateCardJump struct {
	Type     int    `json:"type,omitempty"`
	URL      string `json:"url,omitempty"`
	Title    string `json:"title,omitempty"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}

type messageTemplateCard struct {
	CardType              string                                 `json:"card_type"`
	Source                *messageTemplateCardSource             `json:"source,omitempty"`
	MainTitle             *messageTemplateCardContent            `json:"main_title,omitempty"`
	EmphasisContent       *messageTemplateCardContent            `json:"emphasis_content,omitempty"`
	QuoteArea             *messageTemplateCardQuoteArea          `json:"quote_area,omitempty"`
	SubTitleText          string                                 `json:"sub_title_text,omitempty"`
	HorizontalContentList []messageTemplateCardHorizontalContent `json:"horizontal_content_list,omitempty"`
	JumpList              []messageTemplateCardJump              `json:"jump_list"`
	CardAction            *messageTemplateCardCardAction         `json:"card_action"`
}

type result struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}
