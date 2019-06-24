package namecheap

import (
	"fmt"
	"testing"
)

func TestDomain__GetDomains(t *testing.T) {
	if !clientEnabled {
		t.Skip("namecheap credentials not configured")
	}

	recs, err := testClient.GetDomains()
	for i := range recs {
		fmt.Printf(" %#v\n", recs[i])
	}

	if err != nil {
		t.Fatal(err)
	}
	if len(recs) == 0 {
		t.Fatal("expected domains")
	}
}
