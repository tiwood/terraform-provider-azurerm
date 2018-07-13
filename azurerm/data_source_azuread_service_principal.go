package azurerm

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmActiveDirectoryServicePrincipal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmActiveDirectoryServicePrincipalRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		// TODO: customiseDiff to ensure either `object_id` or `display_name` or `application_id` is set

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"display_name", "application_id"},
			},

			"display_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id", "application_id"},
			},

			"application_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"object_id", "display_name"},
			},
		},
	}
}

func dataSourceArmActiveDirectoryServicePrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).servicePrincipalsClient
	ctx := meta.(*ArmClient).StopContext

	var servicePrincipal *graphrbac.ServicePrincipal

	if v, ok := d.GetOk("object_id"); ok {
		objectId := v.(string)
		app, err := client.Get(ctx, objectId)
		if err != nil {
			if utils.ResponseWasNotFound(app.Response) {
				return fmt.Errorf("Service Principal with Object ID %q was not found!", objectId)
			}

			return fmt.Errorf("Error retrieving Service Principal ID %q: %+v", objectId, err)
		}

		servicePrincipal = &app
	} else {
		apps, err := client.ListComplete(ctx, "")
		if err != nil {
			return fmt.Errorf("Error listing Service Principals: %+v", err)
		}

		if v, ok := d.GetOk("display_name"); ok {
			displayName := v.(string)

			for _, app := range *apps.Response().Value {
				if app.DisplayName == nil {
					continue
				}

				if *app.DisplayName == displayName {
					servicePrincipal = &app
					break
				}
			}

			if servicePrincipal == nil {
				return fmt.Errorf("A Service Principal with the Display Name %q was not found", displayName)
			}
		} else {
			applicationId := d.Get("application_id").(string)

			for _, app := range *apps.Response().Value {
				if app.AppID == nil {
					continue
				}

				if *app.AppID == applicationId {
					servicePrincipal = &app
					break
				}
			}

			if servicePrincipal == nil {
				return fmt.Errorf("A Service Principal for Application ID %q was not found", applicationId)
			}
		}
	}

	d.SetId(*servicePrincipal.ObjectID)

	d.Set("application_id", servicePrincipal.AppID)
	d.Set("display_name", servicePrincipal.DisplayName)
	d.Set("object_id", servicePrincipal.ObjectID)

	return nil
}
