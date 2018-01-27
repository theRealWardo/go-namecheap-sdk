package namecheap

import (
	// "github.com/pearkes/dnsimple/testutil"
	// "strings"
	"testing"
)

func TestHost__GetHosts(t *testing.T) {
	if !clientEnabled {
		t.Skip("namecheap credentials not configured")
	}

	recs, err := testClient.GetHosts(testDomain)
	if err != nil {
		t.Fatal(err)
	}
	if len(recs) != 2 {
		t.Errorf("got %d records", len(recs))
	}
	rec := recs[0]
	testRecord := Record{
		Name:               "www",
		FriendlyName:       "CNAME Record",
		Address:            "parkingpage.namecheap.com.",
		MXPref:             10,
		AssociatedAppTitle: "",
		Id:                 92111926,
		RecordType:         "CNAME",
		TTL:                1800,
		IsActive:           true,
		IsDDNSEnabled:      false,
	}
	if rec.Name != testRecord.Name {
		t.Errorf("%q != %q", rec.Name, testRecord.Name)
	}
	if rec.FriendlyName != testRecord.FriendlyName {
		t.Errorf("%q != %q", rec.FriendlyName, testRecord.FriendlyName)
	}
	if rec.Address != testRecord.Address {
		t.Errorf("%q != %q", rec.Address, testRecord.Address)
	}
	if rec.MXPref != testRecord.MXPref {
		t.Errorf("%q != %q", rec.MXPref, testRecord.MXPref)
	}
	if rec.AssociatedAppTitle != testRecord.AssociatedAppTitle {
		t.Errorf("%q != %q", rec.AssociatedAppTitle, testRecord.AssociatedAppTitle)
	}
	if rec.Id != testRecord.Id {
		t.Errorf("%q != %q", rec.Id, testRecord.Id)
	}
	if rec.RecordType != testRecord.RecordType {
		t.Errorf("%q != %q", rec.RecordType, testRecord.RecordType)
	}
	if rec.TTL != testRecord.TTL {
		t.Errorf("%q != %q", rec.TTL, testRecord.TTL)
	}
	if rec.IsActive != testRecord.IsActive {
		t.Errorf("%v != %v", rec.IsActive, testRecord.IsActive)
	}
	if rec.IsDDNSEnabled != testRecord.IsDDNSEnabled {
		t.Errorf("%v != %v", rec.IsDDNSEnabled, testRecord.IsDDNSEnabled)
	}
}

// func (s *S) Test_SetHosts(c *gocheck.C) {
// 	testServer.Response(200, nil, hostSetExample)
// 	var records []Record

// 	record := Record{
// 		HostName:   "foobar",
// 		RecordType: "CNAME",
// 		Address:    "test.domain.",
// 	}

// 	records = append(records, record)

// 	_, err := s.client.SetHosts("example.com", records)

// 	_ = testServer.WaitRequest()

// 	c.Assert(err, gocheck.IsNil)
// }

// func (s *S) Test_SetHosts_fail(c *gocheck.C) {
// 	testServer.Response(200, nil, hostExampleError)

// 	var records []Record

// 	record := Record{
// 		HostName:   "foobar",
// 		RecordType: "CNAME",
// 		Address:    "test.domain.",
// 	}

// 	records = append(records, record)

// 	_, err := s.client.SetHosts("example.com", records)

// 	_ = testServer.WaitRequest()

// 	c.Assert(strings.Contains(err.Error(), "2019166"), gocheck.Equals, true)
// }

var hostExampleError = `
<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="ERROR" xmlns="http://api.namecheap.com/xml.response">
    <Errors>
        <Error Number="2019166">The domain (huxtest3.com) doesn't seem to be associated with your account.</Error>

	</Errors>
	<Warnings />
	<RequestedCommand>namecheap.domains.dns.setHosts</RequestedCommand>
	<CommandResponse Type="namecheap.domains.dns.setHosts">
		<DomainDNSSetHostsResult Domain="huxtest3.com" EmailType="" IsSuccess="false">
			<Warnings />

		</DomainDNSSetHostsResult>
	</CommandResponse>
	<Server>PHX01SBAPI01</Server>
	<GMTTimeDifference>--5:00</GMTTimeDifference>
	<ExecutionTime>0.025</ExecutionTime>

</ApiResponse>
`

var hostSetExample = `
<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
    <Errors />
    <Warnings />
    <RequestedCommand>namecheap.domains.dns.setHosts</RequestedCommand>
    <CommandResponse Type="namecheap.domains.dns.setHosts">
        <DomainDNSSetHostsResult Domain="example.com" IsSuccess="true">
            <Warnings />

        </DomainDNSSetHostsResult>

    </CommandResponse>
    <Server>PHX01SBAPI01</Server>
    <GMTTimeDifference>--5:00</GMTTimeDifference>
    <ExecutionTime>0.498</ExecutionTime>

</ApiResponse>`

var hostGetExample = `
<?xml version="1.0" encoding="utf-8"?>
<ApiResponse Status="OK" xmlns="http://api.namecheap.com/xml.response">
    <Errors />
    <Warnings />
    <RequestedCommand>namecheap.domains.dns.getHosts</RequestedCommand>
    <CommandResponse Type="namecheap.domains.dns.getHosts">
        <DomainDNSGetHostsResult Domain="huxtest2.com" EmailType="FWD" IsUsingOurDNS="true">
            <host HostId="216107" Name="foobar" Type="CNAME" Address="test.domain." MXPref="10" TTL="1800" AssociatedAppTitle="" FriendlyName="" IsActive="true" IsDDNSEnabled="false" />

        </DomainDNSGetHostsResult>

    </CommandResponse>
    <Server>PHX01SBAPI01</Server>
    <GMTTimeDifference>--5:00</GMTTimeDifference>
    <ExecutionTime>0.704</ExecutionTime>

</ApiResponse>`
