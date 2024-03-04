package domain

import (
	"context"
	"fmt"

	cdnsdk "git.arvancloud.ir/arvancloud/cdn-go-sdk"
)

type DomainService interface {
	Save(DomainEntity) (*cdnsdk.DomainResponse, error)
	FindAll(FindAllDto) (*cdnsdk.DomainsIndex200Response, error)
	Show(domain string) (*cdnsdk.DomainResponse, error)
	Delete(domain string, id string) (*cdnsdk.MessageResponse, error)
}

type domainService struct {
	domains    []DomainEntity
	apiClient  *cdnsdk.APIClient
	domainServ *cdnsdk.DomainStore
}

func NewDomainService(apiClient *cdnsdk.APIClient) DomainService {

	return &domainService{apiClient: apiClient}
}

func (ds *domainService) Save(domain DomainEntity) (*cdnsdk.DomainResponse, error) {
	domainStore := *cdnsdk.NewDomainStore(domain.Domain)
	resp, _, err := ds.apiClient.DomainApi.DomainsStore(context.Background()).DomainStore(domainStore).Execute()
	if err != nil {
		return &cdnsdk.DomainResponse{}, fmt.Errorf("err in calling the cdn sdk: %w", err)
	}
	//ds.domains = append(ds.domains, domain)
	return resp, nil
}

func (ds domainService) FindAll(dto FindAllDto) (*cdnsdk.DomainsIndex200Response, error) {

	var resp *cdnsdk.DomainsIndex200Response
	reqBuilder := ds.apiClient.DomainApi.DomainsIndex(context.Background()).PerPage(dto.PageSize).Page(dto.PageNum)
	if dto.Query != nil {
		reqBuilder = reqBuilder.Search(*dto.Query)
	}
	resp, _, err := reqBuilder.Execute()

	if err != nil {
		return &cdnsdk.DomainsIndex200Response{}, fmt.Errorf("err in calling the cdn sdk: %w", err)
	}
	return resp, nil
}

func (ds domainService) Show(domain string) (*cdnsdk.DomainResponse, error) {
	resp, _, err := ds.apiClient.DomainApi.DomainsShow(context.Background(), domain).Execute()
	if err != nil {
		return &cdnsdk.DomainResponse{}, fmt.Errorf("err in calling the cdn sdk: %w", err)
	}
	return resp, nil
}

func (ds domainService) Delete(domain string, id string) (*cdnsdk.MessageResponse, error) {
	resp, _, err := ds.apiClient.DomainApi.DomainsDestroy(context.Background(), domain).Id(id).Execute()
	if err != nil {
		return &cdnsdk.MessageResponse{}, fmt.Errorf("err in calling the cdn sdk: %w", err)
	}
	return resp, nil
}
