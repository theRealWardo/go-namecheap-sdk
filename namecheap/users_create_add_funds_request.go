package namecheap

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

type CreateAddFundsRequestResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *CreateAddFundsRequestCommandResponse `xml:"CommandResponse"`
}

type CreateAddFundsRequestCommandResponse struct {
	CreateAddFundsRequestResult *CreateAddFundsRequestResult `xml:"Createaddfundsrequestresult"`
}

type CreateAddFundsRequestResult struct {
	TokenID     *string `xml:"TokenID,attr"`
	ReturnURL   *string `xml:"ReturnURL,attr"`
	RedirectURL *string `xml:"RedirectURL,attr"`
}

func (r CreateAddFundsRequestResult) String() string {
	return fmt.Sprintf("{TokenID: %s, ReturnURL: %s, RedirectURL: %s}", *r.TokenID, *r.ReturnURL, *r.RedirectURL)
}

type CreateAddFundsRequestArgs struct {
	PaymentType *string
	Amount      *float64
	ReturnURL   *string
}

func validateCreateAddFundsRequestArgs(args *CreateAddFundsRequestArgs) error {
	if args.PaymentType == nil {
		return fmt.Errorf("PaymentType is required")
	}
	if args.PaymentType != nil && *args.PaymentType != "creditcard" {
		return fmt.Errorf("invalid PaymentType value: %s (only 'creditcard' is allowed)", *args.PaymentType)
	}
	if args.Amount == nil {
		return fmt.Errorf("Amount is required")
	}
	if args.Amount != nil && *args.Amount <= 0 {
		return fmt.Errorf("Amount must be greater than 0")
	}
	if args.ReturnURL == nil {
		return fmt.Errorf("ReturnURL is required")
	}
	return nil
}

func parseCreateAddFundsRequestArgs(args *CreateAddFundsRequestArgs) (*map[string]string, error) {
	params := map[string]string{}

	err := validateCreateAddFundsRequestArgs(args)
	if err != nil {
		return nil, err
	}

	if args.PaymentType != nil {
		params["PaymentType"] = *args.PaymentType
	}
	if args.Amount != nil {
		params["Amount"] = strconv.FormatFloat(*args.Amount, 'f', -1, 64)
	}
	if args.ReturnURL != nil {
		params["ReturnURL"] = *args.ReturnURL
	}

	return &params, nil
}

// CreateAddFundsRequest creates a request to add funds through a credit card
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/users/create-add-funds-request/
func (s *UsersService) CreateAddFundsRequest(args *CreateAddFundsRequestArgs) (*CreateAddFundsRequestCommandResponse, error) {
	var response CreateAddFundsRequestResponse

	params := map[string]string{
		"Command": "namecheap.users.createaddfundsrequest",
	}

	parsedArgs, err := parseCreateAddFundsRequestArgs(args)
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
