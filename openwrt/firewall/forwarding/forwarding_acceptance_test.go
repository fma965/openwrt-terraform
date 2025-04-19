//go:build acceptance.test

package forwarding_test

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
		"src":  lucirpc.String("wan"),
		"dest": lucirpc.String("lan"),
	}
	ok, err := client.CreateSection(ctx, "firwall", "forwarding", "testing", options)
	assert.NilError(t, err)
	assert.Check(t, ok)

	readDataSource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

data "openwrt_firewall_forwarding" "testing" {
	id = "testing"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("data.openwrt_firewall_forwarding.testing", "id", "testing"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_forwarding.testing", "src", "wan"),
			resource.TestCheckResourceAttr("data.openwrt_firewall_forwarding.testing", "dest", "lan"),
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

resource "openwrt_firewall_forwarding" "testing" {
	name = "testing"
	src = "wan"
	dest = "lan"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_forwarding.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_forwarding.testing", "src", "wan"),
			resource.TestCheckResourceAttr("openwrt_firewall_forwarding.testing", "dest", "lan"),
		),
	}
	importValidation := resource.TestStep{
		ImportState:       true,
		ImportStateVerify: true,
		ResourceName:      "openwrt_firewall_forwarding.testing",
	}
	updateAndReadResource := resource.TestStep{
		Config: fmt.Sprintf(`
%s

resource "openwrt_firewall_forwarding" "testing" {	
	name = "testing"
	src = "lan"
	dest = "wan"
}
`,
			providerBlock,
		),
		Check: resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr("openwrt_firewall_forwarding.testing", "name", "testing"),
			resource.TestCheckResourceAttr("openwrt_firewall_forwarding.testing", "src", "lan"),
			resource.TestCheckResourceAttr("openwrt_firewall_forwarding.testing", "dest", "wan"),
		),
	}

	acceptancetest.TerraformSteps(
		t,
		createAndReadResource,
		importValidation,
		updateAndReadResource,
	)
}
