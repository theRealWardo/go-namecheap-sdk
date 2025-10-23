package namecheap

import (
	"encoding/xml"
	"fmt"
)

type GetRegistrarLockResponse struct {
	XMLName         *xml.Name `xml:"ApiResponse"`
	Errors          *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *GetRegistrarLockCommandResponse `xml:"CommandResponse"`
}

type GetRegistrarLockCommandResponse struct {
	Result *GetRegistrarLockResult `xml:"DomainGetRegistrarLockResult"`
}

type GetRegistrarLockResult struct {
	Domain              *string `xml:"Domain,attr"`
	RegistrarLockStatus *bool   `xml:"RegistrarLockStatus,attr"`
}

func (r GetRegistrarLockResult) String() string {
	return fmt.Sprintf("{Domain: %s, RegistrarLockStatus: %t}", *r.Domain, *r.RegistrarLockStatus)
}

// GetRegistrarLock gets the Registrar Lock status for the requested domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/get-registrar-lock/
func (s *DomainsService) GetRegistrarLock(domain string) (*GetRegistrarLockCommandResponse, error) {
	var response GetRegistrarLockResponse

	params := map[string]string{
		"Command":    "namecheap.domains.getRegistrarLock",
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
