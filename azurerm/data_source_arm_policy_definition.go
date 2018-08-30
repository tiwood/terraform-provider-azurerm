package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2016-12-01/policy"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPolicyDefinition() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPolicyDefinitionRead,

		Schema: map[string]*schema.Schema{

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"policy_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmPolicyDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policyDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	var policyDefinition policy.Definition

	name := d.Get("name").(string)
	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: ARM policy definition with name %q was not found", name)
		}

		return fmt.Errorf("Error making Read request on ARM policy definition with name %q: %+v", name, err)
	}

	policyDefinition = resp

	d.SetId(*policyDefinition.ID)
	d.Set("name", policyDefinition.Name)
	d.Set("display_name", policyDefinition.DisplayName)
	d.Set("policy_type", policyDefinition.PolicyType)

	return nil
}
