package namecheap

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

// SetEmailForwardingResponse represents the API response for setEmailForwarding
type SetEmailForwardingResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *SetEmailForwardingCommandResponse `xml:"CommandResponse"`
}

// SetEmailForwardingCommandResponse wraps the result
type SetEmailForwardingCommandResponse struct {
	DomainDNSSetEmailForwardingResult *SetEmailForwardingResult `xml:"DomainDNSSetEmailForwardingResult"`
}

// SetEmailForwardingResult contains the result of setting email forwarding
type SetEmailForwardingResult struct {
	Domain    *string `xml:"Domain,attr"`
	IsSuccess *bool   `xml:"IsSuccess,attr"`
}

func (r SetEmailForwardingResult) String() string {
	return fmt.Sprintf("{Domain: %s, IsSuccess: %t}", *r.Domain, *r.IsSuccess)
}

// EmailForwardingEntry represents a single mailbox to forward-to mapping
type EmailForwardingEntry struct {
	Mailbox   string
	ForwardTo string
}

// SetEmailForwarding sets email forwarding for a domain name.
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains-dns/set-email-forwarding/
func (dds *DomainsDNSService) SetEmailForwarding(domainName string, forwardingRules []EmailForwardingEntry) (*SetEmailForwardingCommandResponse, error) {
	var response SetEmailForwardingResponse

	params := map[string]string{
		"Command":    "namecheap.domains.dns.setEmailForwarding",
		"DomainName": domainName,
	}

	for i, rule := range forwardingRules {
		index := strconv.Itoa(i + 1)
		params["mailbox"+index] = rule.Mailbox
		params["ForwardTo"+index] = rule.ForwardTo
	}

	_, err := dds.client.DoXML(params, &response)
	if err != nil {
		return nil, err
	}
	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil
}
