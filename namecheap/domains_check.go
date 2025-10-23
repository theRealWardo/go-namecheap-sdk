package namecheap

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type CheckResponse struct {
	XMLName         *xml.Name `xml:"ApiResponse"`
	Errors          *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *CheckCommandResponse `xml:"CommandResponse"`
}

type CheckCommandResponse struct {
	DomainCheckResults *[]DomainCheckResult `xml:"DomainCheckResult"`
}

type DomainCheckResult struct {
	Domain                    *string `xml:"Domain,attr"`
	Available                 *string `xml:"Available,attr"`
	ErrorNo                   *string `xml:"ErrorNo,attr"`
	Description               *string `xml:"Description,attr"`
	IsPremiumName             *string `xml:"IsPremiumName,attr"`
	PremiumRegistrationPrice  *string `xml:"PremiumRegistrationPrice,attr"`
	PremiumRenewalPrice       *string `xml:"PremiumRenewalPrice,attr"`
	PremiumRestorePrice       *string `xml:"PremiumRestorePrice,attr"`
	PremiumTransferPrice      *string `xml:"PremiumTransferPrice,attr"`
	IcannFee                  *string `xml:"IcannFee,attr"`
	EapFee                    *string `xml:"EapFee,attr"`
}

func (r DomainCheckResult) String() string {
	domain := ""
	if r.Domain != nil {
		domain = *r.Domain
	}
	available := ""
	if r.Available != nil {
		available = *r.Available
	}
	isPremium := ""
	if r.IsPremiumName != nil {
		isPremium = *r.IsPremiumName
	}
	return fmt.Sprintf("{Domain: %s, Available: %s, IsPremiumName: %s}", domain, available, isPremium)
}

// Check checks the availability of one or more domains
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/check/
func (s *DomainsService) Check(domains []string) (*CheckCommandResponse, error) {
	var response CheckResponse

	params := map[string]string{
		"Command":    "namecheap.domains.check",
		"DomainList": strings.Join(domains, ","),
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
