package namecheap

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type ReactivateArgs struct {
	PromotionCode   *string
	YearsToAdd      *int
	IsPremiumDomain *bool
	PremiumPrice    *string
}

type ReactivateResponse struct {
	XMLName         *xml.Name `xml:"ApiResponse"`
	Errors          *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *ReactivateCommandResponse `xml:"CommandResponse"`
}

type ReactivateCommandResponse struct {
	DomainReactivateResult *ReactivateResult `xml:"DomainReactivateResult"`
}

type ReactivateResult struct {
	Domain        *string `xml:"Domain,attr"`
	IsSuccess     *bool   `xml:"IsSuccess,attr"`
	ChargedAmount *string `xml:"ChargedAmount,attr"`
	OrderID       *int    `xml:"OrderID,attr"`
	TransactionID *int    `xml:"TransactionID,attr"`
}

func (r ReactivateResult) String() string {
	domain := ""
	if r.Domain != nil {
		domain = *r.Domain
	}
	success := false
	if r.IsSuccess != nil {
		success = *r.IsSuccess
	}
	amount := ""
	if r.ChargedAmount != nil {
		amount = *r.ChargedAmount
	}
	return fmt.Sprintf("{Domain: %s, IsSuccess: %t, ChargedAmount: %s}", domain, success, amount)
}

func validateReactivateArgs(args *ReactivateArgs) error {
	if args.IsPremiumDomain != nil && *args.IsPremiumDomain {
		if args.PremiumPrice == nil {
			return fmt.Errorf("PremiumPrice is required when IsPremiumDomain is true")
		}
	}
	return nil
}

func parseReactivateArgs(args *ReactivateArgs) (*map[string]string, error) {
	params := map[string]string{}

	err := validateReactivateArgs(args)
	if err != nil {
		return nil, err
	}

	if args.PromotionCode != nil {
		params["PromotionCode"] = *args.PromotionCode
	}

	if args.YearsToAdd != nil {
		params["YearsToAdd"] = strconv.Itoa(*args.YearsToAdd)
	}

	if args.IsPremiumDomain != nil {
		params["IsPremiumDomain"] = strconv.FormatBool(*args.IsPremiumDomain)
	}

	if args.PremiumPrice != nil {
		params["PremiumPrice"] = *args.PremiumPrice
	}

	return &params, nil
}

// Reactivate reactivates an expired domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/reactivate/
func (s *DomainsService) Reactivate(domain string, args *ReactivateArgs) (*ReactivateCommandResponse, error) {
	var response ReactivateResponse

	params := map[string]string{
		"Command": "namecheap.domains.reactivate",
	}

	params["DomainName"] = domain

	if args != nil {
		parsedArgs, err := parseReactivateArgs(args)
		if err != nil {
			return nil, err
		}
		for k, v := range *parsedArgs {
			params[k] = v
		}
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
