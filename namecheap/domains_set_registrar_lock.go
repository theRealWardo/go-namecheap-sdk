package namecheap

import (
	"encoding/xml"
	"fmt"
)

type LockAction string

const (
	LockActionLock   LockAction = "LOCK"
	LockActionUnlock LockAction = "UNLOCK"
)

type SetRegistrarLockResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *SetRegistrarLockCommandResponse `xml:"CommandResponse"`
}

type SetRegistrarLockCommandResponse struct {
	Result *SetRegistrarLockResult `xml:"DomainSetRegistrarLockResult"`
}

type SetRegistrarLockResult struct {
	Domain    *string `xml:"Domain,attr"`
	IsSuccess *bool   `xml:"IsSuccess,attr"`
}

func (r SetRegistrarLockResult) String() string {
	return fmt.Sprintf("{Domain: %s, IsSuccess: %t}", *r.Domain, *r.IsSuccess)
}

// SetRegistrarLock sets the Registrar Lock status for a domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/set-registrar-lock/
func (s *DomainsService) SetRegistrarLock(domain string, lockAction *LockAction) (*SetRegistrarLockCommandResponse, error) {
	var response SetRegistrarLockResponse

	params := map[string]string{
		"Command":    "namecheap.domains.setRegistrarLock",
		"DomainName": domain,
	}

	if lockAction != nil {
		params["LockAction"] = string(*lockAction)
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
