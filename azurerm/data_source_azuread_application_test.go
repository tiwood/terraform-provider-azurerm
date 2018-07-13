package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMAzureADApplication_byObjectId(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADApplication_objectId(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_basic(id),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("http://acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADApplication_byObjectIdComplete(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADApplication_objectIdComplete(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_basic(id),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("http://homepage-%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMAzureADApplication_byName(t *testing.T) {
	dataSourceName := "data.azurerm_azuread_application.test"
	id := uuid.New().String()
	config := testAccDataSourceAzureRMAzureADApplication_name(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMActiveDirectoryApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMActiveDirectoryApplication_basic(id),
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMActiveDirectoryApplicationExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", fmt.Sprintf("acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "homepage", fmt.Sprintf("http://acctest%s", id)),
					resource.TestCheckResourceAttr(dataSourceName, "identifier_uris.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "reply_urls.#", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "oauth2_allow_implicit_flow", "false"),
					resource.TestCheckResourceAttrSet(dataSourceName, "application_id"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMAzureADApplication_objectId(id string) string {
	template := testAccAzureRMActiveDirectoryApplication_basic(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_application" "test" {
  object_id = "${azurerm_azuread_application.test.id}"
}
`, template)
}

func testAccDataSourceAzureRMAzureADApplication_objectIdComplete(id string) string {
	template := testAccAzureRMActiveDirectoryApplication_complete(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_application" "test" {
  object_id = "${azurerm_azuread_application.test.id}"
}
`, template)
}

func testAccDataSourceAzureRMAzureADApplication_name(id string) string {
	template := testAccAzureRMActiveDirectoryApplication_basic(id)
	return fmt.Sprintf(`
%s

data "azurerm_azuread_application" "test" {
  name = "${azurerm_azuread_application.test.name}"
}
`, template)
}
