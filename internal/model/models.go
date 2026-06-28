package model

type GetLinkShortenReq struct {
	Url string `json:"url"`
}

type GetLinkShortenRes struct {
	Result string `json:"result"`
}

type Link struct {
	Id          string `json:"uuid"`
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}
