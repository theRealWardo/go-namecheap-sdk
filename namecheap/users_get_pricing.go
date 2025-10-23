package namecheap

import (
	"encoding/xml"
	"fmt"
)

type ProductType string

const (
	ProductTypeDomain         ProductType = "DOMAIN"
	ProductTypeSSLCertificate ProductType = "SSLCERTIFICATE"
)

type ProductCategory string

const (
	ProductCategoryDomains ProductCategory = "DOMAINS"
	ProductCategoryComodo  ProductCategory = "COMODO"
)

type ActionName string

const (
	ActionNameRegister   ActionName = "REGISTER"
	ActionNameRenew      ActionName = "RENEW"
	ActionNameReactivate ActionName = "REACTIVATE"
	ActionNameTransfer   ActionName = "TRANSFER"
	ActionNamePurchase   ActionName = "PURCHASE"
)

type ProductName string

const (
	ProductNameCom        ProductName = "COM"
	ProductNameInstantSSL ProductName = "INSTANTSSL"
)

type GetPricingArgs struct {
	ProductType     ProductType
	ProductCategory *ProductCategory
	PromotionCode   *string
	ActionName      *ActionName
	ProductName     *ProductName
}

type GetPricingResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *GetPricingCommandResponse `xml:"CommandResponse"`
}

type GetPricingCommandResponse struct {
	UserGetPricingResult *GetPricingResult `xml:"UserGetPricingResult"`
}

type GetPricingResult struct {
	ProductTypes *[]ProductTypeResult `xml:"ProductType"`
}

type ProductTypeResult struct {
	Name              *string                  `xml:"Name,attr"`
	ProductCategories *[]ProductCategoryResult `xml:"ProductCategory"`
}

type ProductCategoryResult struct {
	Name     *string          `xml:"Name,attr"`
	Products *[]ProductResult `xml:"Product"`
}

type ProductResult struct {
	Name   *string        `xml:"Name,attr"`
	Prices *[]PriceResult `xml:"Price"`
}

type PriceResult struct {
	Duration     *string `xml:"Duration,attr"`
	DurationType *string `xml:"DurationType,attr"`
	Price        *string `xml:"Price,attr"`
	RegularPrice *string `xml:"RegularPrice,attr"`
	YourPrice    *string `xml:"YourPrice,attr"`
	CouponPrice  *string `xml:"CouponPrice,attr"`
	Currency     *string `xml:"Currency,attr"`
}

func (r GetPricingResult) String() string {
	if r.ProductTypes == nil || len(*r.ProductTypes) == 0 {
		return "{ProductTypes: []}"
	}
	return fmt.Sprintf("{ProductTypes: %d types}", len(*r.ProductTypes))
}

func validateGetPricingArgs(args *GetPricingArgs) error {
	if args.ProductType == "" {
		return fmt.Errorf("ProductType is required")
	}
	return nil
}

func parseGetPricingArgs(args *GetPricingArgs) (*map[string]string, error) {
	params := map[string]string{}

	err := validateGetPricingArgs(args)
	if err != nil {
		return nil, err
	}

	params["ProductType"] = string(args.ProductType)

	if args.ProductCategory != nil {
		params["ProductCategory"] = string(*args.ProductCategory)
	}

	if args.PromotionCode != nil {
		params["PromotionCode"] = *args.PromotionCode
	}

	if args.ActionName != nil {
		params["ActionName"] = string(*args.ActionName)
	}

	if args.ProductName != nil {
		params["ProductName"] = string(*args.ProductName)
	}

	return &params, nil
}

// GetPricing returns pricing information for a requested product type
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/users/get-pricing/
func (s *UsersService) GetPricing(args *GetPricingArgs) (*GetPricingCommandResponse, error) {
	var response GetPricingResponse

	params := map[string]string{
		"Command": "namecheap.users.getPricing",
	}

	parsedArgs, err := parseGetPricingArgs(args)
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
