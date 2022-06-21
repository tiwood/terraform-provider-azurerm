package network_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type IPGroupCidrResource struct{}

func TestAccIpGroupCidr_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group_cidr", "test")
	r := IPGroupCidrResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIpGroupCidr_multiple(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group_cidr", "test")
	r := IPGroupCidrResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.multiple(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group.test").Key("tags.env").HasValue("prod"),
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
				check.That("azurerm_ip_group_cidr.multiple_1").ExistsInAzure(r),
				check.That("azurerm_ip_group_cidr.multiple_2").ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIpGroupCidr_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_ip_group_cidr", "test")
	r := IPGroupCidrResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That("azurerm_ip_group_cidr.test").ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_ip_group_cidr"),
		},
	})
}

func (t IPGroupCidrResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.IpGroupCidrID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Network.IPGroupsClient.Get(ctx, id.ResourceGroup, id.IpGroupName, "")
	if err != nil {
		return nil, fmt.Errorf("reading IP Group (%s): %+v", id, err)
	}

	if !utils.SliceContainsValue(*resp.IPAddresses, state.Attributes["cidr"]) {
		return utils.Bool(false), nil
	}

	return utils.Bool(true), nil
}

func (IPGroupCidrResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-network-%d"
  location = "%s"
}

resource "azurerm_ip_group" "test" {
  name                = "acceptanceTestIpGroup1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tags = {
    env = "prod"
  }
}

resource "azurerm_ip_group_cidr" "test" {
  ip_group_id = azurerm_ip_group.test.id
  cidr        = "10.0.0.0/24"
}


`, data.RandomInteger, data.Locations.Primary)
}

func (r IPGroupCidrResource) multiple(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group_cidr" "multiple_1" {
  ip_group_id = azurerm_ip_group.test.id
  cidr        = "10.10.0.0/24"
}

resource "azurerm_ip_group_cidr" "multiple_2" {
  ip_group_id = azurerm_ip_group.test.id
  cidr        = "10.20.0.0/24"
}
`, r.basic(data))
}

func (r IPGroupCidrResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_ip_group_cidr" "import" {
  ip_group_id = azurerm_ip_group_cidr.test.ip_group_id
  cidr        = azurerm_ip_group_cidr.test.cidr
}
`, r.basic(data))
}
