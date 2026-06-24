package model

type GetLinkShortenReq struct {
	Url string `json:"url"`
}

type GetLinkShortenRes struct {
	Result string `json:"result"`
}
