package namecheap

import (
	"encoding/xml"
	"fmt"
)

type GetBalancesResponse struct {
	XMLName         *xml.Name `xml:"ApiResponse"`
	Errors          *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *GetBalancesCommandResponse `xml:"CommandResponse"`
}

type GetBalancesCommandResponse struct {
	UserGetBalancesResult *GetBalancesResult `xml:"UserGetBalancesResult"`
}

type GetBalancesResult struct {
	Currency                  *string `xml:"Currency,attr"`
	AvailableBalance          *string `xml:"AvailableBalance,attr"`
	AccountBalance            *string `xml:"AccountBalance,attr"`
	EarnedAmount              *string `xml:"EarnedAmount,attr"`
	WithdrawableAmount        *string `xml:"WithdrawableAmount,attr"`
	FundsRequiredForAutoRenew *string `xml:"FundsRequiredForAutoRenew,attr"`
}

func (r GetBalancesResult) String() string {
	return fmt.Sprintf("{Currency: %s, AvailableBalance: %s, AccountBalance: %s, EarnedAmount: %s, WithdrawableAmount: %s, FundsRequiredForAutoRenew: %s}",
		*r.Currency, *r.AvailableBalance, *r.AccountBalance, *r.EarnedAmount, *r.WithdrawableAmount, *r.FundsRequiredForAutoRenew)
}

// GetBalances gets information about fund in the user's account
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/users/get-balances/
func (s *UsersService) GetBalances() (*GetBalancesCommandResponse, error) {
	var response GetBalancesResponse

	params := map[string]string{
		"Command": "namecheap.users.getBalances",
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
