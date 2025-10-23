package namecheap

import (
	"encoding/xml"
	"fmt"
)

// GetEmailForwardingResponse represents the API response for getEmailForwarding
type GetEmailForwardingResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *GetEmailForwardingCommandResponse `xml:"CommandResponse"`
}

// GetEmailForwardingCommandResponse wraps the result
type GetEmailForwardingCommandResponse struct {
	DomainDNSGetEmailForwardingResult *GetEmailForwardingResult `xml:"DomainDNSGetEmailForwardingResult"`
}

// GetEmailForwardingResult contains the email forwarding details
type GetEmailForwardingResult struct {
	Domain   *string                `xml:"Domain,attr"`
	Forwards *[]EmailForwardingRule `xml:"Forward"`
}

// EmailForwardingRule represents a single email forwarding rule
type EmailForwardingRule struct {
	Mailbox     *string `xml:"mailbox,attr"`
	ForwardTo   *string `xml:",chardata"`
}

// GetEmailForwarding retrieves email forwarding settings for the requested domain.
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains-dns/get-email-forwarding/
func (dds *DomainsDNSService) GetEmailForwarding(domainName string) (*GetEmailForwardingCommandResponse, error) {
	var response GetEmailForwardingResponse

	params := map[string]string{
		"Command":    "namecheap.domains.dns.getEmailForwarding",
		"DomainName": domainName,
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
