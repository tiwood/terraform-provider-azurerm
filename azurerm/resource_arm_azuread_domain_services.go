package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2017-01-01/aad"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDomainServices() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDomainServicesCreate,
		Read:   resourceArmDomainServicesRead,
		Delete: resourceArmDomainServicesDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"virtual_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"internal_ip_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceArmDomainServicesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).domainServicesClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	domainName := d.Get("domain_name").(string)
	tenantID := d.Get("tenant_id").(string)
	virtualNetworkID := d.Get("virtual_network_id").(string)
	subnetID := d.Get("subnet_id").(string)

	//create the properties object
	properties := aad.DomainServiceProperties{
		TenantID:   &tenantID,
		DomainName: &domainName,
		VnetSiteID: &virtualNetworkID,
		SubnetID:   &subnetID,
	}

	//create the Azure AD Domain Services resource (long-running operation)
	future, err := client.CreateOrUpdate(ctx, resourceGroupName, domainName, properties)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating AAD Domain Services %q (Resource Group %q): %+v", domainName, resourceGroupName, err)
	}

	//the creation of Azure AD Domain Services is a long-running operation, we're going to wait for completion
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of AAD Domain Services %q (Resource Group %q): %+v", domainName, resourceGroupName, err)
	}

	//the creation should be complete, now we're going to get the resource and set the objectID
	read, err := client.Get(ctx, resourceGroupName, domainName)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Error reading the objectID from AAD Domain Services %q (Resource Group %q): %+v", domainName, resourceGroupName, err)
	}

	//set the resource id
	d.SetId(*read.ID)

	return resourceArmDomainServicesRead(d, meta)
}

func resourceArmDomainServicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).domainServicesClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	domainName := d.Get("domain_name").(string)

	read, err := client.Get(ctx, resourceGroupName, domainName)
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] Azure AD Domain Services with name %q (Resource Group %q) was not found - removing from state", domainName, resourceGroupName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Azure AD Domain Services with name %q (Resource Group %q): %+v", domainName, resourceGroupName, err)
	}

	//update the schema with the retrieved properties
	d.Set("domain_name", read.DomainName)
	d.Set("tenant_id", read.TenantID)
	d.Set("virtual_network_id", read.VnetSiteID)
	d.Set("subnet_id", read.SubnetID)

	internalIPAddresses := make([]string, 0)
	if s := read.DomainControllerIPAddress; s != nil {
		internalIPAddresses = *s
	}
	if err := d.Set("internal_ip_address", internalIPAddresses); err != nil {
		return fmt.Errorf("Error setting `internal_ip_address`: %+v", err)
	}

	return nil
}

func resourceArmDomainServicesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).domainServicesClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroupName := d.Get("resource_group_name").(string)
	domainName := d.Get("domain_name").(string)

	//delete the Azure AD Domain Services resource (long-running operation)
	future, err := client.Delete(ctx, resourceGroupName, domainName)
	if err != nil {
		return fmt.Errorf("Error Deleting AAD Domain Services %q (Resource Group %q): %+v", domainName, resourceGroupName, err)
	}

	//the deletion of Azure AD Domain Services is a long-running operation, we're going to wait for completion
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for the deletion of AAD Domain Services %q (Resource Group %q): %+v", domainName, resourceGroupName, err)
	}

	return nil
}
