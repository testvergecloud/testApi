package server

import (
	cdnsdk "git.arvancloud.ir/arvancloud/cdn-go-sdk"
	"go-starter/config"
)

func NewCdnApiClient(config *config.Config) *cdnsdk.APIClient {
	configuration := cdnsdk.NewConfiguration()
	configuration.AddDefaultHeader("Authorization", config.CDNApiKey)
	return cdnsdk.NewAPIClient(configuration)
}
