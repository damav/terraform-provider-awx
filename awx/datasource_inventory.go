package awx

import (
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "gitlab.com/dhendel/awx-go"
)

func dataSourceInventory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceInventoryRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this inventory",
			},
			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of the ansible inventory",
			},
		},
	}
}

func dataSourceInventoryRead(d *schema.ResourceData, meta interface{}) error {
	awx := meta.(*awxgo.AWX)
	awxService := awx.InventoriesService
	_, res, err := awxService.ListInventories(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d.SetId(strconv.Itoa(res.Results[0].ID))
	d = setInventoryDataSource(d, res.Results[0])
	return nil
}

func setInventoryDataSource(d *schema.ResourceData, r *awxgo.Inventory) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("id", r.ID)
	return d
}
