package namecheap

import (
	"fmt"
	"strconv"
)

type ContactInfo struct {
	FirstName           *string
	LastName            *string
	Address1            *string
	Address2            *string
	City                *string
	StateProvince       *string
	StateProvinceChoice *string
	PostalCode          *string
	Country             *string
	Phone               *string
	PhoneExt            *string
	Fax                 *string
	EmailAddress        *string
	OrganizationName    *string
	JobTitle            *string
}

type CreateArgs struct {
	DomainName    *string
	Years         *int
	PromotionCode *string

	Registrant *ContactInfo
	Tech       *ContactInfo
	Admin      *ContactInfo
	AuxBilling *ContactInfo

	AddFreeWhoisguard *bool
	WGEnabled         *bool

	Nameservers *string

	IdnCode *string

	IsPremiumDomain *bool
	PremiumPrice    *string
	EapFee          *string
}

type DomainsCreateResponse struct {
	Errors *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsCreateCommandResponse `xml:"CommandResponse"`
}

type DomainsCreateCommandResponse struct {
	DomainCreateResult *DomainsCreateResult `xml:"DomainCreateResult"`
}

type DomainsCreateResult struct {
	Domain            *string `xml:"Domain,attr"`
	Registered        *bool   `xml:"Registered,attr"`
	ChargedAmount     *string `xml:"ChargedAmount,attr"`
	DomainID          *int    `xml:"DomainID,attr"`
	OrderID           *int    `xml:"OrderID,attr"`
	TransactionID     *int    `xml:"TransactionID,attr"`
	WhoisguardEnable  *bool   `xml:"WhoisguardEnable,attr"`
	NonRealTimeDomain *bool   `xml:"NonRealTimeDomain,attr"`
}

func (r DomainsCreateResult) String() string {
	domain := ""
	if r.Domain != nil {
		domain = *r.Domain
	}
	registered := false
	if r.Registered != nil {
		registered = *r.Registered
	}
	chargedAmount := ""
	if r.ChargedAmount != nil {
		chargedAmount = *r.ChargedAmount
	}
	return fmt.Sprintf("{Domain: %s, Registered: %t, ChargedAmount: %s}", domain, registered, chargedAmount)
}

func validateContactInfo(contact *ContactInfo, prefix string) error {
	if contact == nil {
		return fmt.Errorf("%s contact information is required", prefix)
	}
	if contact.FirstName == nil || *contact.FirstName == "" {
		return fmt.Errorf("%sFirstName is required", prefix)
	}
	if contact.LastName == nil || *contact.LastName == "" {
		return fmt.Errorf("%sLastName is required", prefix)
	}
	if contact.Address1 == nil || *contact.Address1 == "" {
		return fmt.Errorf("%sAddress1 is required", prefix)
	}
	if contact.City == nil || *contact.City == "" {
		return fmt.Errorf("%sCity is required", prefix)
	}
	if contact.StateProvince == nil || *contact.StateProvince == "" {
		return fmt.Errorf("%sStateProvince is required", prefix)
	}
	if contact.PostalCode == nil || *contact.PostalCode == "" {
		return fmt.Errorf("%sPostalCode is required", prefix)
	}
	if contact.Country == nil || *contact.Country == "" {
		return fmt.Errorf("%sCountry is required", prefix)
	}
	if contact.Phone == nil || *contact.Phone == "" {
		return fmt.Errorf("%sPhone is required", prefix)
	}
	if contact.EmailAddress == nil || *contact.EmailAddress == "" {
		return fmt.Errorf("%sEmailAddress is required", prefix)
	}
	return nil
}

func addContactToParams(params map[string]string, contact *ContactInfo, prefix string) {
	if contact.FirstName != nil {
		params[prefix+"FirstName"] = *contact.FirstName
	}
	if contact.LastName != nil {
		params[prefix+"LastName"] = *contact.LastName
	}
	if contact.Address1 != nil {
		params[prefix+"Address1"] = *contact.Address1
	}
	if contact.Address2 != nil {
		params[prefix+"Address2"] = *contact.Address2
	}
	if contact.City != nil {
		params[prefix+"City"] = *contact.City
	}
	if contact.StateProvince != nil {
		params[prefix+"StateProvince"] = *contact.StateProvince
	}
	if contact.StateProvinceChoice != nil {
		params[prefix+"StateProvinceChoice"] = *contact.StateProvinceChoice
	}
	if contact.PostalCode != nil {
		params[prefix+"PostalCode"] = *contact.PostalCode
	}
	if contact.Country != nil {
		params[prefix+"Country"] = *contact.Country
	}
	if contact.Phone != nil {
		params[prefix+"Phone"] = *contact.Phone
	}
	if contact.PhoneExt != nil {
		params[prefix+"PhoneExt"] = *contact.PhoneExt
	}
	if contact.Fax != nil {
		params[prefix+"Fax"] = *contact.Fax
	}
	if contact.EmailAddress != nil {
		params[prefix+"EmailAddress"] = *contact.EmailAddress
	}
	if contact.OrganizationName != nil {
		params[prefix+"OrganizationName"] = *contact.OrganizationName
	}
	if contact.JobTitle != nil {
		params[prefix+"JobTitle"] = *contact.JobTitle
	}
}

func validateCreateArgs(args *CreateArgs) error {
	if args.DomainName == nil || *args.DomainName == "" {
		return fmt.Errorf("DomainName is required")
	}
	if args.Years == nil {
		return fmt.Errorf("Years is required")
	}
	if *args.Years < 1 {
		return fmt.Errorf("Years must be at least 1")
	}

	if err := validateContactInfo(args.Registrant, "Registrant"); err != nil {
		return err
	}
	if err := validateContactInfo(args.Tech, "Tech"); err != nil {
		return err
	}
	if err := validateContactInfo(args.Admin, "Admin"); err != nil {
		return err
	}
	if err := validateContactInfo(args.AuxBilling, "AuxBilling"); err != nil {
		return err
	}

	return nil
}

func parseCreateArgs(args *CreateArgs) (*map[string]string, error) {
	params := map[string]string{}

	err := validateCreateArgs(args)
	if err != nil {
		return nil, err
	}

	params["DomainName"] = *args.DomainName
	params["Years"] = strconv.Itoa(*args.Years)

	if args.PromotionCode != nil {
		params["PromotionCode"] = *args.PromotionCode
	}

	addContactToParams(params, args.Registrant, "Registrant")
	addContactToParams(params, args.Tech, "Tech")
	addContactToParams(params, args.Admin, "Admin")
	addContactToParams(params, args.AuxBilling, "AuxBilling")

	if args.AddFreeWhoisguard != nil {
		if *args.AddFreeWhoisguard {
			params["AddFreeWhoisguard"] = "yes"
		} else {
			params["AddFreeWhoisguard"] = "no"
		}
	}

	if args.WGEnabled != nil {
		if *args.WGEnabled {
			params["WGEnabled"] = "yes"
		} else {
			params["WGEnabled"] = "no"
		}
	}

	if args.Nameservers != nil {
		params["Nameservers"] = *args.Nameservers
	}

	if args.IdnCode != nil {
		params["IdnCode"] = *args.IdnCode
	}

	if args.IsPremiumDomain != nil {
		params["IsPremiumDomain"] = strconv.FormatBool(*args.IsPremiumDomain)
	}

	if args.PremiumPrice != nil {
		params["PremiumPrice"] = *args.PremiumPrice
	}

	if args.EapFee != nil {
		params["EapFee"] = *args.EapFee
	}

	return &params, nil
}

// Create registers a new domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains/create/
func (s *DomainsService) Create(args *CreateArgs) (*DomainsCreateCommandResponse, error) {
	var response DomainsCreateResponse

	params := map[string]string{
		"Command": "namecheap.domains.create",
	}

	parsedArgs, err := parseCreateArgs(args)
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
