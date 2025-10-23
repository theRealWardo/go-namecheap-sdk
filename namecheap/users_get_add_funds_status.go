package namecheap

import (
	"encoding/xml"
	"fmt"
)

type GetAddFundsStatusResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *GetAddFundsStatusCommandResponse `xml:"CommandResponse"`
}

type GetAddFundsStatusCommandResponse struct {
	GetAddFundsStatusResult *GetAddFundsStatusResult `xml:"GetAddFundsStatusResult"`
}

type GetAddFundsStatusResult struct {
	TransactionID *string `xml:"TransactionID,attr"`
	Amount        *string `xml:"Amount,attr"`
	Status        *string `xml:"Status,attr"`
}

func (r GetAddFundsStatusResult) String() string {
	return fmt.Sprintf("{TransactionID: %s, Amount: %s, Status: %s}", *r.TransactionID, *r.Amount, *r.Status)
}

// GetAddFundsStatus gets the status of add funds request
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/users/get-add-funds-status/
func (s *UsersService) GetAddFundsStatus(tokenID string) (*GetAddFundsStatusCommandResponse, error) {
	var response GetAddFundsStatusResponse

	params := map[string]string{
		"Command": "namecheap.users.getAddFundsStatus",
		"TokenID": tokenID,
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
