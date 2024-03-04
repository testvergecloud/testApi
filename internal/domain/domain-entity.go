package domain

type DomainEntity struct {
	Domain     string `json:"domain"`
	DomainType string `json:"domain_type"`
}

func NewDomain(domain string, domainType string) DomainEntity {
	return DomainEntity{Domain: domain, DomainType: domainType}
}
