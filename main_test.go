package main

import (
	"context"
	"os"
	"testing"
	"time"

	logf "github.com/cert-manager/cert-manager/pkg/logs"
	acmetest "github.com/cert-manager/cert-manager/test/acme"
	dns "github.com/cert-manager/cert-manager/test/acme"
	"github.com/cert-manager/cert-manager/test/acme/server"

	testDNS "github.com/pmarques/cert-manager-webhook-namecheap/test"
)

const (
	chanllengeKey = "123abc=="
)

var (
	zone = os.Getenv("TEST_ZONE_NAME")
	fqdn = "dns01-test." + zone
)

func TestRunsSuite(t *testing.T) {
	ctx := logf.NewContext(context.TODO(), logf.Log, t.Name())

	dnsServer := &server.BasicServer{
		Handler: &testDNS.Handler{
			Log: logf.FromContext(ctx, "testDNSServer"),
			TxtRecords: map[string][][]string{
				fqdn: {
					{},
					{},
					{},
					{chanllengeKey},
					{chanllengeKey},
				},
			},
		},
		Zones: []string{zone},
	}

	if err := dnsServer.Run(ctx, "udp"); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}

	defer dnsServer.Shutdown()

	fixture := acmetest.NewFixture(&namecheapDNSProviderSolver{},
		acmetest.SetResolvedZone(zone),
		acmetest.SetResolvedFQDN(fqdn),
		acmetest.SetDNSChallengeKey(chanllengeKey),
		acmetest.SetManifestPath("testdata/namecheap"),
		acmetest.SetDNSServer(dnsServer.ListenAddr()),
		dns.SetPropagationLimit(time.Duration(10)*time.Second),
		acmetest.SetUseAuthoritative(false),
	)

	fixture.RunConformance(t)
}
