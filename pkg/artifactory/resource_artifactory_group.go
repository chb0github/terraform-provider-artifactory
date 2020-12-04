package artifactory

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"log"
)

func resourceArtifactoryGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreateOrUpdate,
		Read:   resourceGroupRead,
		Update: resourceGroupCreateOrUpdate,
		Delete: resourceGroupDelete,
		Exists: resourceGroupExists,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_join": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"admin_privileges": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"realm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateLowerCase,
			},
			"realm_attributes": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func unmarshalGroup(s *schema.ResourceData) (*services.Group, error) {
	d := &ResourceData{s}

	group := services.Group{
		Name:            d.getStringRef("name", false),
		Description:     d.getStringRef("description", false),
		AutoJoin:        d.getBoolRef("auto_join", false),
		AdminPrivileges: d.getBoolRef("admin_privileges", false),
		Realm:           d.getStringRef("realm", false),
		RealmAttributes: d.getStringRef("realm_attributes", false),
	}
	// Validator
	if group.AdminPrivileges != nil && group.AutoJoin != nil &&
		*group.AdminPrivileges && *group.AutoJoin {
		return nil, fmt.Errorf("error: auto_join cannot be true if admin_privileges is true")
	}

	return &group, nil
}

func resourceGroupCreateOrUpdate(d *schema.ResourceData, m interface{}) error {

	client := *m.(*ArtClient).ArtNew

	group, err := unmarshalGroup(d)

	if err != nil {
		return err
	}

	if err := client.CreateGroup(*group); err != nil {
		return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {

			if err := client.CreateGroup(*group); err != nil {
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(resourceGroupRead(d, m))
		})
	}
	d.SetId(*group.Name)
	return nil
}

func resourceGroupRead(d *schema.ResourceData, m interface{}) error {
	client := *m.(*ArtClient).ArtNew
	group, err := client.GetGroup(d.Id())

	if err != nil {
		return err
	}

	// If we 404 it is likely the resources was externally deleted
	// If the ID is updated to blank, this tells Terraform the resource no longer exist
	if group == nil {
		d.SetId("")
		return fmt.Errorf("no group returned, it may have been externally deleted")
	}

	hasErr := false
	logError := cascadingErr(&hasErr)
	logError(d.Set("name", group.Name))
	logError(d.Set("description", group.Description))
	logError(d.Set("auto_join", group.AutoJoin))
	logError(d.Set("admin_privileges", group.AdminPrivileges))
	logError(d.Set("realm", group.Realm))
	logError(d.Set("realm_attributes", group.RealmAttributes))
	if hasErr {
		return fmt.Errorf("failed to marshal group")
	}
	return nil
}

func resourceGroupDelete(group *schema.ResourceData, m interface{}) error {
	client := *m.(*ArtClient).ArtNew
	log.Printf("Deleting group %s", group.Id())
	return client.DeleteGroup(group.Id())
}

func resourceGroupExists(group *schema.ResourceData, m interface{}) (bool, error) {
	client := *m.(*ArtClient).ArtNew
	log.Printf("Check group %s", group.Id())
	return client.GroupExists(group.Id())
}
