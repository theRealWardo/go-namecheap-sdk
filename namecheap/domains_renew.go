package namecheap

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type RenewArgs struct {
	Years            *int
	PromotionCode    *string
	IsPremiumDomain  *bool
	PremiumPrice     *string
}

type RenewResponse struct {
	XMLName         *xml.Name `xml:"ApiResponse"`
	Errors          *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *RenewCommandResponse `xml:"CommandResponse"`
}

type RenewCommandResponse struct {
	DomainRenewResult *RenewResult `xml:"DomainRenewResult"`
}

type RenewResult struct {
	DomainName    *string        `xml:"DomainName,attr"`
	DomainID      *int           `xml:"DomainID,attr"`
	Renew         *bool          `xml:"Renew,attr"`
	OrderID       *int           `xml:"OrderID,attr"`
	TransactionID *int           `xml:"TransactionID,attr"`
	ChargedAmount *string        `xml:"ChargedAmount,attr"`
	DomainDetails *DomainDetails `xml:"DomainDetails"`
}

type DomainDetails struct {
	ExpiredDate *string `xml:"ExpiredDate"`
	NumYears    *int    `xml:"NumYears"`
}

func (r RenewResult) String() string {
	return fmt.Sprintf("{DomainName: %s, DomainID: %d, Renew: %t, OrderID: %d, TransactionID: %d, ChargedAmount: %s}",
		*r.DomainName, *r.DomainID, *r.Renew, *r.OrderID, *r.TransactionID, *r.ChargedAmount)
}

func validateRenewArgs(args *RenewArgs) error {
	if args.Years == nil {
		return fmt.Errorf("Years is required")
	}
	if *args.Years < 1 || *args.Years > 10 {
		return fmt.Errorf("Years must be between 1 and 10")
	}
	if args.IsPremiumDomain != nil && *args.IsPremiumDomain && args.PremiumPrice == nil {
		return fmt.Errorf("PremiumPrice is required when IsPremiumDomain is true")
	}
	return nil
}

func parseRenewArgs(args *RenewArgs) (*map[string]string, error) {
	params := map[string]string{}

	err := validateRenewArgs(args)
	if err != nil {
		return nil, err
	}

	if args.Years != nil {
		params["Years"] = strconv.Itoa(*args.Years)
	}

	if args.PromotionCode != nil {
		params["PromotionCode"] = *args.PromotionCode
	}

	if args.IsPremiumDomain != nil {
		params["IsPremiumDomain"] = strconv.FormatBool(*args.IsPremiumDomain)
	}

	if args.PremiumPrice != nil {
		params["PremiumPrice"] = *args.PremiumPrice
	}

	return &params, nil
}

// Renew renews an expiring domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/renew/
func (s *DomainsService) Renew(domain string, args *RenewArgs) (*RenewCommandResponse, error) {
	var response RenewResponse

	params := map[string]string{
		"Command": "namecheap.domains.renew",
	}

	params["DomainName"] = domain

	parsedArgs, err := parseRenewArgs(args)
	if err != nil {
		return nil, err
	}
	for k, v := range *parsedArgs {
		params[k] = v
	}

	_, err = s.client.DoXML(params, &response)
	if err != nil {
		return nil, err
	}

	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil
}
