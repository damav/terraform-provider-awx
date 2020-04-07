package awx

import (
	"fmt"
	"strconv"
	"time"

	awxgo "github.com/davidfischer-ch/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProjectObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Delete: resourceProjectDelete,
		Update: resourceProjectUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this project",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this project.",
			},
			"local_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Local path (relative to PROJECTS_ROOT) containing playbooks and related files for this project.",
			},
			"scm_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the source control system used to store the project (one of '', git, hg, svn, insights).",
			},
			"scm_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The location where the project is stored.",
			},
			"scm_branch": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Specific branch, tag or commit to checkout.",
			},
			"scm_refspec": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "For git projects, an additional refspec to fetch.",
			},
			"scm_clean": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Discard any local changes before syncing the project.",
			},
			"scm_delete_on_update": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Delete the project before syncing.",
			},
			"credential_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Numeric ID of the project credential",
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The amount of time (in seconds) to run before the task is canceled.",
			},
			"organization_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Numeric ID of the project organization",
			},
			"scm_update_on_launch": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Update the project when a job is launched that uses the project.",
			},
			"scm_update_cache_timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The number of seconds after the last project update ran that a new project update will be launched as a job dependency.",
			},
			"allow_override": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Allow changing the SCM branch or revision in a job template that uses this project.",
			},
			"custom_virtualenv": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     nil,
				Description: "Local absolute file path containing a custom Python virtualenv to use",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.ProjectService

	_, res, err := awxService.ListProjects(map[string]string{
		"name":         d.Get("name").(string),
		"organization": strconv.Itoa(d.Get("organization_id").(int))},
	)
	if err != nil {
		return err
	}
	if len(res.Results) >= 1 {
		return fmt.Errorf("Project with name %s already exists in the organization %d",
			d.Get("name").(string), d.Get("organization_id").(int))
	}

	result, err := awxService.CreateProject(map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"local_path":               d.Get("local_path").(string),
		"scm_type":                 d.Get("scm_type").(string),
		"scm_url":                  d.Get("scm_url").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"scm_refspec":              d.Get("scm_refspec").(string),
		"scm_clean":                d.Get("scm_clean").(bool),
		"scm_delete_on_update":     d.Get("scm_delete_on_update").(bool),
		"credential":               d.Get("credential_id").(int),
		"timeout":                  d.Get("timeout").(int),
		"organization":             d.Get("organization_id").(int),
		"scm_update_on_launch":     d.Get("scm_update_on_launch").(bool),
		"scm_update_cache_timeout": d.Get("scm_update_cache_timeout").(int),
		"allow_override":           d.Get("allow_override").(bool),
		"custom_virtualenv":        d.Get("custom_virtualenv").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceProjectRead(d, m)
}

func resourceProjectUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.ProjectService
	_, res, err := awxService.ListProjects(map[string]string{
		"id":           d.Id(),
		"organization": strconv.Itoa(d.Get("organization_id").(int))},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("Project with name %s doesn't exists in the organization %s",
			d.Get("name").(string), d.Get("organization_id").(string))
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.UpdateProject(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"local_path":               d.Get("local_path").(string),
		"scm_type":                 d.Get("scm_type").(string),
		"scm_url":                  d.Get("scm_url").(string),
		"scm_branch":               d.Get("scm_branch").(string),
		"scm_refspec":              d.Get("scm_refspec").(string),
		"scm_clean":                d.Get("scm_clean").(bool),
		"scm_delete_on_update":     d.Get("scm_delete_on_update").(bool),
		"credential":               d.Get("credential_id").(int),
		"timeout":                  d.Get("timeout").(int),
		"organization":             d.Get("organization_id").(int),
		"scm_update_on_launch":     d.Get("scm_update_on_launch").(bool),
		"scm_update_cache_timeout": d.Get("scm_update_cache_timeout").(int),
		"allow_override":           d.Get("allow_override").(bool),
		"custom_virtualenv":        d.Get("custom_virtualenv").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	return resourceProjectRead(d, m)
}

func resourceProjectRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.ProjectService
	_, res, err := awxService.ListProjects(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setProjectResourceData(d, res.Results[0])
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.ProjectService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	var jobID int
	var finished time.Time
	_, res, err := awxService.ListProjects(map[string]string{
		"name": d.Get("name").(string),
		"id":   d.Id()},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		d.SetId("")
		return nil
	}
	if res.Results[0].SummaryFields.CurrentJob["id"] != nil {
		jobID = int(res.Results[0].SummaryFields.CurrentJob["id"].(float64))
	} else if res.Results[0].SummaryFields.LastJob["id"] != nil {
		jobID = int(res.Results[0].SummaryFields.LastJob["id"].(float64))
	}
	if jobID != 0 {
		_, err = awx.ProjectUpdatesService.ProjectUpdateCancel(jobID)
		if err != nil {
			return err
		}
	}
	// check if finished is 0
	for finished.IsZero() {
		prj, _ := awx.ProjectUpdatesService.ProjectUpdateGet(jobID)
		finished = prj.Finished
		time.Sleep(1 * time.Second)
	}

	if _, err = awxService.DeleteProject(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setProjectResourceData(d *schema.ResourceData, r *awxgo.Project) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("local_path", r.LocalPath)
	d.Set("scm_type", r.ScmType)
	d.Set("scm_url", r.ScmURL)
	d.Set("scm_branch", r.ScmBranch)
	d.Set("scm_refspec", r.ScmRefSpec)
	d.Set("scm_clean", r.ScmClean)
	d.Set("scm_delete_on_update", r.ScmDeleteOnUpdate)
	d.Set("credential_id", r.Credential)
	d.Set("timeout", r.Timeout)
	d.Set("organization_id", r.Organization)
	d.Set("scm_update_on_launch", r.ScmUpdateOnLaunch)
	d.Set("scm_update_cache_timeout", r.ScmUpdateCacheTimeout)
	d.Set("allow_override", r.AllowOverride)
	d.Set("custom_virtualenv", r.CustomVirtualenv)
	return d
}
