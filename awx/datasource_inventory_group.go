package awx

import (
	"strconv"

	awxgo "github.com/davidfischer-ch/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceInventoryGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInventoryGroupRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this group",
			},
			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of the ansible inventory group",
			},
			"inventory_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of the ansible inventory this group belongs to",
			},
		},
	}
}

func dataSourceInventoryGroupRead(d *schema.ResourceData, meta interface{}) error {
	awx := meta.(*awxgo.AWX)
	awxService := awx.GroupService
	_, res, err := awxService.ListGroups(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}

	d.SetId(strconv.Itoa(res.Results[0].ID))
	d = setInventoryGroupSourceData(d, res.Results[0])
	return nil
}

func setInventoryGroupSourceData(d *schema.ResourceData, r *awxgo.Group) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("id", r.ID)
	d.Set("inventory_id", r.Inventory)
	return d
}
