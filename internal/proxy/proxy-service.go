package proxy

import (
	"fmt"
)

type ProxyService interface {
	Redirect() (Entity, error)
	TargetUrl() string
}

type proxyService struct {
	urlTarget string
	apiToken  string
}

func NewProxyService(urlTarget string, apiToken string) ProxyService {
	return &proxyService{urlTarget: urlTarget, apiToken: apiToken}
}

func (ps *proxyService) Redirect() (Entity, error) {
	urlTarget := ps.urlTarget
	apiToken := ps.apiToken

	fmt.Println(urlTarget, apiToken)

	proxy := Entity{urlTarget, apiToken}
	return proxy, nil
}

func (ps *proxyService) TargetUrl() string {
	return ps.urlTarget
}
