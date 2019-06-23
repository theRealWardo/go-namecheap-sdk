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
	fmt.Printf("%v", recs)
	fmt.Printf("%v", testClient)

	if err != nil {
		t.Fatal(err)
	}
	if len(recs) == 0 {
		t.Fatal("expected domains")
	}
}
