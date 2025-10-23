package namecheap

import (
	"encoding/xml"
	"fmt"
)

type GetTldListResponse struct {
	XMLName         *xml.Name `xml:"ApiResponse"`
	Errors          *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *GetTldListCommandResponse `xml:"CommandResponse"`
}

type GetTldListCommandResponse struct {
	Tlds *GetTldListResult `xml:"Tlds"`
}

type GetTldListResult struct {
	Tlds *[]Tld `xml:"Tld"`
}

type Tld struct {
	Name                           *string `xml:"Name,attr"`
	NonRealTime                    *bool   `xml:"NonRealTime,attr"`
	MinRegisterYears               *int    `xml:"MinRegisterYears,attr"`
	MaxRegisterYears               *int    `xml:"MaxRegisterYears,attr"`
	MinRenewYears                  *int    `xml:"MinRenewYears,attr"`
	MaxRenewYears                  *int    `xml:"MaxRenewYears,attr"`
	MinTransferYears               *int    `xml:"MinTransferYears,attr"`
	MaxTransferYears               *int    `xml:"MaxTransferYears,attr"`
	IsApiRegisterable              *bool   `xml:"IsApiRegisterable,attr"`
	IsApiRenewable                 *bool   `xml:"IsApiRenewable,attr"`
	IsApiTransferable              *bool   `xml:"IsApiTransferable,attr"`
	IsEppRequired                  *bool   `xml:"IsEppRequired,attr"`
	IsDisableModContact            *bool   `xml:"IsDisableModContact,attr"`
	IsDisableWGAllot               *bool   `xml:"IsDisableWGAllot,attr"`
	IsIncludeInExtendedSearchOnly  *bool   `xml:"IsIncludeInExtendedSearchOnly,attr"`
	SequenceNumber                 *int    `xml:"SequenceNumber,attr"`
	Type                           *string `xml:"Type,attr"`
	IsSupportsIDN                  *bool   `xml:"IsSupportsIDN,attr"`
	Category                       *string `xml:"Category,attr"`
	Description                    *string `xml:",chardata"`
}

func (r GetTldListResult) String() string {
	if r.Tlds == nil {
		return "{Tlds: []}"
	}
	return fmt.Sprintf("{Tlds: %d items}", len(*r.Tlds))
}

// GetTldList returns a list of TLDs
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/get-tld-list/
func (s *DomainsService) GetTldList() (*GetTldListCommandResponse, error) {
	var response GetTldListResponse

	params := map[string]string{
		"Command": "namecheap.domains.getTldList",
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
