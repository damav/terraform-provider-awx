package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_job_template test case
func TestAccAWXJobTemplate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccJobTemplateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateJobTemplate("name", "alpha"),
					testAccCheckStateJobTemplate("description", "Alpha job template example"),
					testAccCheckStateJobTemplate("job_type", "run"),
					testAccCheckStateJobTemplate("inventory_id", "1"),
					testAccCheckStateJobTemplate("project_id", "4"),
					testAccCheckStateJobTemplate("playbook", "hello_world.yml"),
				),
			},
		},
	})
}

func testAccCheckStateJobTemplate(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_job_template.alpha"]
		if !ok {
			return fmt.Errorf("awx_job_template.alpha not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cr := rs.Primary

		if cr.Attributes[skey] != svalue {
			return fmt.Errorf("%s != %s (actual: %s)", skey, svalue, cr.Attributes[skey])
		}

		return nil
	}
}

const testAccJobTemplateConfig = `
resource "awx_job_template" "alpha" {
	name         = "alpha"
	description  = "Alpha job template example"
	project_id   = 4
	job_type     = "run"
	inventory_id = 1
	playbook     = "hello_world.yml"
  }
`
