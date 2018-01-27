package namecheap

import (
	"encoding/xml"
)

type RecordsResponse struct {
	XMLName xml.Name `xml:"ApiResponse"`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse struct {
		Records []Record `xml:"DomainDNSGetHostsResult>host"`
	} `xml:"CommandResponse"`
}

type RecordsCreateResult struct {
	XMLName xml.Name `xml:"ApiResponse"`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse struct {
		DomainDNSSetHostsResult struct {
			Domain    string `xml:"Domain,attr"`
			IsSuccess bool   `xml:"IsSuccess,attr"`
		} `xml:"DomainDNSSetHostsResult"`
	} `xml:"CommandResponse"`
}

type NSListResponse struct {
	XMLName xml.Name `xml:"ApiResponse`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse struct {
		DomainDNSGetListResult []string `xml:"DomainDNSGetListResult>Nameserver"`
	} `sml:CommandResponse"`
}

type NSSetCustomRepsonse struct {
	XMLName xml.Name `xml:"ApiResponse`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse struct {
		DomainDNSSetCustomResult struct {
			Domain  string `xml:"Domain,attr"`
			Updated bool   `xml:"Updated,attr"`
		} `xml:"DomainDNSSetCustomResult"`
	} `xml:"CommandResponse"`
}

type NSSetDefaultResponse struct {
	XMLName xml.Name `xml:"ApiResponse`
	Errors  []struct {
		Message string `xml:",chardata"`
		Number  string `xml:"Number,attr"`
	} `xml:"Errors>Error"`
	CommandResponse struct {
		DomainDNSSetDefaultResult struct {
			Domain  string `xml:"Domain,attr"`
			Updated bool   `xml:"Updated,attr"`
		} `xml:"DomainDNSSetDefaultResult"`
	} `xml:"CommandResponse"`
}

// Record is used to represent a retrieved Record. All properties
// are set as strings.
type Record struct {
	Name               string `xml:"Name,attr"`
	FriendlyName       string `xml:"FriendlyName,attr"`
	Address            string `xml:"Address,attr"`
	MXPref             int    `xml:"MXPref,attr"`
	AssociatedAppTitle string `xml:"AssociatedAppTitle,attr"`
	Id                 int    `xml:"HostId,attr"`
	RecordType         string `xml:"Type,attr"`
	TTL                int    `xml:"TTL,attr"`
	IsActive           bool   `xml:"IsActive,attr"`
	IsDDNSEnabled      bool   `xml:"IsDDNSEnabled,attr"`
}
