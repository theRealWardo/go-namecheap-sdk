package namecheap

import (
	"encoding/xml"
	"fmt"
)

type DomainsGetContactsResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsGetContactsCommandResponse `xml:"CommandResponse"`
}

type DomainsGetContactsCommandResponse struct {
	DomainContactsResult *DomainsGetContactsResult `xml:"DomainContactsResult"`
}

type DomainsGetContactsResult struct {
	Domain            *string                `xml:"Domain,attr"`
	DomainNameID      *string                `xml:"domainnameid,attr"`
	Registrant        *DomainContactInfo     `xml:"Registrant"`
	Tech              *DomainContactInfo     `xml:"Tech"`
	Admin             *DomainContactInfo     `xml:"Admin"`
	AuxBilling        *DomainContactInfo     `xml:"AuxBilling"`
	CurrentAttributes *CurrentAttributes     `xml:"CurrentAttributes"`
	WhoisGuardContact *WhoisGuardContactInfo `xml:"WhoisGuardContact"`
}

type DomainContactInfo struct {
	ReadOnly            *string `xml:"ReadOnly,attr"`
	OrganizationName    *string `xml:"OrganizationName"`
	JobTitle            *string `xml:"JobTitle"`
	FirstName           *string `xml:"FirstName"`
	LastName            *string `xml:"LastName"`
	Address1            *string `xml:"Address1"`
	Address2            *string `xml:"Address2"`
	City                *string `xml:"City"`
	StateProvince       *string `xml:"StateProvince"`
	StateProvinceChoice *string `xml:"StateProvinceChoice"`
	PostalCode          *string `xml:"PostalCode"`
	Country             *string `xml:"Country"`
	Phone               *string `xml:"Phone"`
	Fax                 *string `xml:"Fax"`
	EmailAddress        *string `xml:"EmailAddress"`
	PhoneExt            *string `xml:"PhoneExt"`
}

type CurrentAttributes struct {
	RegistrantNexus        *string `xml:"RegistrantNexus"`
	RegistrantNexusCountry *string `xml:"RegistrantNexusCountry"`
	RegistrantPurpose      *string `xml:"RegistrantPurpose"`
}

type WhoisGuardContactInfo struct {
	Registrant        *DomainContactInfo `xml:"Registrant"`
	Tech              *DomainContactInfo `xml:"Tech"`
	Admin             *DomainContactInfo `xml:"Admin"`
	AuxBilling        *DomainContactInfo `xml:"AuxBilling"`
	CurrentAttributes *CurrentAttributes `xml:"CurrentAttributes"`
}

func (r DomainsGetContactsResult) String() string {
	return fmt.Sprintf("{Domain: %s, Registrant: %v, Tech: %v, Admin: %v, AuxBilling: %v}",
		stringValue(r.Domain), r.Registrant, r.Tech, r.Admin, r.AuxBilling)
}

func (c DomainContactInfo) String() string {
	return fmt.Sprintf("{Name: %s %s, Email: %s, Organization: %s}",
		stringValue(c.FirstName), stringValue(c.LastName), stringValue(c.EmailAddress), stringValue(c.OrganizationName))
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// GetContacts gets contact information for the requested domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/get-contacts/
func (s *DomainsService) GetContacts(domain string) (*DomainsGetContactsCommandResponse, error) {
	var response DomainsGetContactsResponse

	params := map[string]string{
		"Command":    "namecheap.domains.getContacts",
		"DomainName": domain,
	}

	_, err := s.client.DoXML(params, &response)
	if err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil
}
