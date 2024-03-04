package proxy

type Entity struct {
	UrlTarget string `json:"url_target"`
	ApiToken  string `json:"api_token"`
}

func NewDomain(urlTarget string, apiToken string) Entity {
	return Entity{UrlTarget: urlTarget, ApiToken: apiToken}
}
