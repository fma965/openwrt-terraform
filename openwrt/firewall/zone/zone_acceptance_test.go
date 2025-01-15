//go:build acceptance.test

package zone_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"gotest.tools/v3/assert"
)

func TestDataSourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	client := openWrtServer.LuCIRPCClient(
		ctx,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()
	options := lucirpc.Options{
		"name":    lucirpc.String("testing"),
		"forward": lucirpc.String("ACCEPT"),
		"input":   lucirpc.String("ACCEPT"),
		"output":  lucirpc.String("ACCEPT"),
	}
	ok, err := client.CreateSection(ctx, "firwall", "zone", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_firewall_zone" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_firewall_zone.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_zone.testing", "name", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_zone.testing", "forward", "ACCEPT"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_zone.testing", "input", "ACCEPT"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_zone.testing", "output", "ACCEPT"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		readDataSource,
	)
}

func TestResourceAcceptance(t *testing.T) {
	ctx := context.Background()
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		*dockerPool,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()

	createAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_zone" "testing" {
	name = "testing"
	forward = "ACCEPT"
	input = "ACCEPT"
	output = "ACCEPT"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "forward", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "input", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "output", "ACCEPT"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_firewall_zone.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_zone" "testing" {
	name = "testing"
	forward = "ACCEPT"
	input = "ACCEPT"
	output = "ACCEPT"
	network = [
		"vlan0",
		"vlan1",
	]
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "id", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "forward", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "input", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "output", "ACCEPT"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "network[0]", "vlan0"),
			resource.TestCheckResourceAttr("openwrt_firewall_zone.testing", "network[1]", "vlan1"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
