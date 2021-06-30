package namecheap

import (
	"encoding/xml"
	"fmt"
	"strconv"
)

const (
	minTTL int = 60
	maxTTL int = 60000
)

var allowedRecordTypes = []string{"A", "AAAA", "ALIAS", "CAA", "CNAME", "MX", "MXE", "NS", "TXT", "URL", "URL301", "FRAME"}
var allowedEmailTypeValues = []string{"NONE", "MXE", "MX", "FWD", "OX"}
var allowedTagValues = []string{"issue", "issuewild", "iodef"}

type DomainsDNSHostRecord struct {
	// Sub-domain/hostname to create the record for
	HostName *string
	// Possible values: A, AAAA, ALIAS, CAA, CNAME, MX, MXE, NS, TXT, URL, URL301, FRAME
	RecordType *string
	// Possible values are URL or ClientIp address. The value for this parameter is based on RecordType.
	Address *string
	// MX preference for host. Applicable for MX records only.
	MXPref *uint8
	// Time to live for all record types.Possible values: any value between 60 to 60000
	// Default Value: 1800 (if 0 value has been provided)
	TTL *int
}

type DomainsDNSSetHostsArgs struct {
	// Domain to setHosts
	Domain *string
	// DomainsDNSHostRecord list
	Records *[]DomainsDNSHostRecord
	// Possible values are MXE, MX, FWD, OX or NONE
	// If empty, then this field won't be forwarded
	EmailType *string
	// Is an unsigned integer between 0 and 255.
	// The flag value is an 8-bit number, the most significant bit of which indicates the criticality of understanding of a record by a CA.
	// It's recommended to use '0'
	// If nil provided, then this field is ignored
	Flag *uint8
	// A non-zero sequence of US-ASCII letters and numbers in lower case. The tag value can be one of the following values:
	// "issue" — specifies the certification authority that is authorized to issue a certificate for the domain name or subdomain record used in the title.
	// "issuewild" — specifies the certification authority that is allowed to issue a wildcard certificate for the domain name or subdomain record used in the title. The certificate applies to the domain name or subdomain directly and to all its subdomains.
	// "iodef" — specifies the e-mail address or URL (compliant with RFC 5070) a CA should use to notify a client if any issuance policy violation spotted by this CA.
	Tag *string
}

type DomainsDNSSetHostsResponse struct {
	XMLName *xml.Name `xml:"ApiResponse"`
	Errors  *[]struct {
		Message *string `xml:",chardata"`
		Number  *string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse *DomainsDNSSetHostsCommandResponse `xml:"CommandResponse"`
}

type DomainsDNSSetHostsCommandResponse struct {
	DomainDNSSetHostsResult *DomainDNSSetHostsResult `xml:"DomainDNSSetHostsResult"`
}

type DomainDNSSetHostsResult struct {
	Domain *string `xml:"Domain,attr"`
	//EmailType *string `xml:"EmailType,attr"`
	IsSuccess *bool `xml:"IsSuccess,attr"`
}

func (d DomainDNSSetHostsResult) String() string {
	return fmt.Sprintf("{Domain: %s, IsSuccess: %t}", *d.Domain, *d.IsSuccess)
}

// SetHosts sets DNS host records settings for the requested domain
//
// Namecheap doc: https://www.namecheap.com/support/api/methods/domains-dns/set-hosts/
func (dds DomainsDNSService) SetHosts(args *DomainsDNSSetHostsArgs) (*DomainsDNSSetHostsCommandResponse, error) {
	var response DomainsDNSSetHostsResponse

	params := map[string]string{
		"Command": "namecheap.domains.dns.setHosts",
	}

	// parse input arguments
	parsedArgsMap, err := parseDomainsDNSSetHostsArgs(args)
	if err != nil {
		return nil, err
	}

	// merge parsed arguments with params
	for k, v := range *parsedArgsMap {
		params[k] = v
	}

	req, err := dds.client.NewRequest(params)
	if err != nil {
		return nil, err
	}
	resp, err := dds.client.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = decodeBody(resp.Body, &response)
	if err != nil {
		return nil, err
	}
	if response.Errors != nil && len(*response.Errors) > 0 {
		apiErr := (*response.Errors)[0]
		return nil, fmt.Errorf("%s (%s)", *apiErr.Message, *apiErr.Number)
	}

	return response.CommandResponse, nil
}

func parseDomainsDNSSetHostsArgs(args *DomainsDNSSetHostsArgs) (*map[string]string, error) {
	params := map[string]string{}

	parsedDomain, err := ParseDomain(*args.Domain)
	if err != nil {
		return nil, err
	}

	params["SLD"] = parsedDomain.SLD
	params["TLD"] = parsedDomain.TLD

	if args.EmailType != nil {
		if isValidEmailType(*args.EmailType) {
			params["EmailType"] = *args.EmailType
		} else {
			return nil, fmt.Errorf("invalid EmailType value: %s", *args.EmailType)
		}
	}

	if args.Flag != nil {
		params["Flag"] = strconv.Itoa(int(*args.Flag))
	}

	if args.Tag != nil {
		if isValidTagValue(*args.Tag) {
			params["Tag"] = *args.Tag
		} else {
			return nil, fmt.Errorf("invalid Tag value: %s", *args.Tag)
		}
	}

	if args.Records != nil {
		for i, record := range *args.Records {
			recordIndexString := strconv.Itoa(i + 1)

			if record.HostName != nil {
				params["HostName"+recordIndexString] = *record.HostName
			} else {
				return nil, fmt.Errorf("Records[%d].HostName is required", i)
			}

			if record.RecordType == nil {
				return nil, fmt.Errorf("Records[%d].RecordType is required", i)
			}

			if isValidRecordType(*record.RecordType) {
				params["RecordType"+recordIndexString] = *record.RecordType
			} else {
				return nil, fmt.Errorf("invalid Records[%d].RecordType value: %s", i, *record.RecordType)
			}

			if record.TTL != nil {
				if *record.TTL >= minTTL && *record.TTL <= maxTTL {
					params["TTL"+recordIndexString] = strconv.Itoa(*record.TTL)
				} else {
					return nil, fmt.Errorf("invalid Records[%d].TTL value: %d", i, *record.TTL)
				}
			}

			if record.Address != nil {
				params["Address"+recordIndexString] = *record.Address
			} else {
				return nil, fmt.Errorf("Records[%d].Address is required", i)
			}

			if record.MXPref != nil {
				params["MXPref"+recordIndexString] = strconv.Itoa(int(*record.MXPref))
			}

		}
	}

	return &params, nil
}

func isValidEmailType(emailType string) bool {
	for _, value := range allowedEmailTypeValues {
		if emailType == value {
			return true
		}
	}
	return false
}

func isValidTagValue(tag string) bool {
	for _, value := range allowedTagValues {
		if tag == value {
			return true
		}
	}
	return false
}

func isValidRecordType(recordType string) bool {
	for _, value := range allowedRecordTypes {
		if recordType == value {
			return true
		}
	}
	return false
}
