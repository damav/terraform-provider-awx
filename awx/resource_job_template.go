package awx

import (
	"fmt"
	"strconv"
	"time"

	awxgo "github.com/davidfischer-ch/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceJobTemplateObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateCreate,
		Read:   resourceJobTemplateRead,
		Delete: resourceJobTemplateDelete,
		Update: resourceJobTemplateUpdate,
		Importer: &schema.ResourceImporter{
			State: importJobTemplateData,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this job template.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Optional description of this job template.",
			},
			"job_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of: run, check, scan",
			},
			"inventory_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ask_inventory_on_launch"},
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"playbook": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"scm_branch": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Branch to use in job run. Project default used if blank. Only allowed if project allow_override field is set to true.",
			},
			"forks": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"verbosity": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "In range 0-5 (Normal, Verbose, More Verbose, Debug, Connection Debug, WinRM Debug)",
			},
			"extra_vars": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"job_tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"force_handlers": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"skip_tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"start_at_task": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"timeout": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The amount of time (in seconds) to run before the task is canceled.",
			},
			"use_fact_cache": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, Tower will act as an Ansible Fact Cache Plugin; persisting facts at the end of a playbook run to the database and caching facts for use by Ansible.",
			},
			"host_config_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_scm_branch_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_diff_mode_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_variables_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_skip_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_job_type_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_verbosity_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_inventory_on_launch": &schema.Schema{
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"inventory_id"},
			},
			"ask_credential_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"survey_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"become_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"diff_mode": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, textual changes made to any templated files on the host are shown in the standard output.",
			},
			"allow_simultaneous": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"custom_virtualenv": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Local absolute file path containing a custom Python virtualenv to use.",
			},
			"job_slice_count": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The number of jobs to slice into at runtime. Will cause the Job Template to launch a workflow if value is greater than 1.",
			},
			"webhook_service": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Service that webhook requests will be accepted from (github or gitlab)",
			},
			"webhook_credential_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Personal Access Token for posting back the status to the service API.",
			},

			// Extra fields (such as self identifier)
			"job_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"extra_credential_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"vault_credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceJobTemplateCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	var jobID int
	var finished time.Time
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"name":    d.Get("name").(string),
		"project": d.Get("project_id").(string)},
	)

	if err != nil {
		return err
	}

	if len(res.Results) >= 1 {
		return fmt.Errorf("JobTemplate with name %s already exists",
			d.Get("name").(string))
	}
	_, prj, err := awx.ProjectService.ListProjects(map[string]string{
		"id": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if prj.Results[0].SummaryFields.CurrentJob["id"] != nil {
		jobID = int(prj.Results[0].SummaryFields.CurrentJob["id"].(float64))
	} else if prj.Results[0].SummaryFields.LastJob["id"] != nil {
		jobID = int(prj.Results[0].SummaryFields.LastJob["id"].(float64))
	}

	if jobID != 0 {
		// check if finished is 0
		for finished.IsZero() {
			prj, _ := awx.ProjectUpdatesService.ProjectUpdateGet(jobID)
			if prj != nil {
				finished = prj.Finished
				time.Sleep(1 * time.Second)
			}
		}
	}

	payload := map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  AtoipOr(d.Get("project_id").(string), nil),
		"playbook":                 d.Get("playbook").(string),
		"credential":               AtoipOr(d.Get("credential_id").(string), nil),
		"scm_branch":               d.Get("scm_branch").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_scm_branch_on_launch": d.Get("ask_scm_branch_on_launch").(bool),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"job_slice_count":          d.Get("job_slice_count").(int),
		"webhook_service":          d.Get("webhook_service").(string),
		"webhook_credential":       AtoipOr(d.Get("webhook_credential_id").(string), nil),
		"vault_credential":         AtoipOr(d.Get("vault_credential_id").(string), nil),
	}

	result, err := awxService.CreateJobTemplate(payload, map[string]string{})
	if err != nil {
		return err
	}

	if creds, ok := d.GetOk("extra_credential_ids"); ok {
		for _, c := range creds.([]interface{}) {
			_, err := awxService.AddJobTemplateCredential(result.ID, c.(int))
			if err != nil {
				return err
			}
		}

	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id":      d.Id(),
		"project": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("JobTemplate with name %s doesn't exists",
			d.Get("name").(string))
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	result, err := awxService.UpdateJobTemplate(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  AtoipOr(d.Get("project_id").(string), nil),
		"playbook":                 d.Get("playbook").(string),
		"credential":               AtoipOr(d.Get("credential_id").(string), nil),
		"scm_branch":               d.Get("scm_branch").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_scm_branch_on_launch": d.Get("ask_scm_branch_on_launch").(bool),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"job_slice_count":          d.Get("job_slice_count").(int),
		"webhook_service":          d.Get("webhook_service").(string),
		"webhook_credential":       AtoipOr(d.Get("webhook_credential_id").(string), nil),
		"vault_credential":         AtoipOr(d.Get("vault_credential_id").(string), nil),
	}, map[string]string{})
	if err != nil {
		return err
	}

	if creds, ok := d.GetOk("extra_credential_ids"); ok {
		for _, c := range creds.([]interface{}) {
			_, err := awxService.AddJobTemplateCredential(result.ID, c.(int))
			if err != nil {
				return err
			}
		}

	}

	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id": strconv.Itoa(d.Get("job_id").(int)),
		//"name":    d.Get("name").(string),
		//"project": d.Get("project_id").(string),
	})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setJobTemplateResourceData(d, res.Results[0])
	return nil
}

func resourceJobTemplateDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id":      d.Id(),
		"project": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		d.SetId("")
		return nil
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.DeleteJobTemplate(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setJobTemplateResourceData(d *schema.ResourceData, r *awxgo.JobTemplate) *schema.ResourceData {
	d.Set("job_id", r.ID)
	d.Set("allow_simultaneous", r.AllowSimultaneous)
	d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch)
	d.Set("ask_diff_mode_on_launch", r.AskDiffModeOnLaunch)
	d.Set("ask_inventory_on_launch", r.AskInventoryOnLaunch)
	d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch)
	d.Set("ask_limit_on_launch", r.AskLimitOnLaunch)
	d.Set("ask_scm_branch_on_launch", r.AskScmBranchOnLaunch)
	d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch)
	d.Set("ask_tags_on_launch", r.AskTagsOnLaunch)
	d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch)
	d.Set("ask_verbosity_on_launch", r.AskVerbosityOnLaunch)
	d.Set("become_enabled", r.BecomeEnabled)
	d.Set("become_enabled", r.BecomeEnabled)
	d.Set("credential_id", r.Credential)
	d.Set("custom_virtualenv", r.CustomVirtualenv)
	d.Set("description", r.Description)
	d.Set("diff_mode", r.DiffMode)
	d.Set("diff_mode", r.DiffMode)
	d.Set("extra_vars", r.ExtraVars)
	d.Set("force_handlers", r.ForceHandlers)
	d.Set("forks", r.Forks)
	d.Set("host_config_key", r.HostConfigKey)
	d.Set("inventory_id", r.Inventory)
	d.Set("job_slice_count", r.JobSliceCount)
	d.Set("job_tags", r.JobTags)
	d.Set("job_type", r.JobType)
	d.Set("limit", r.Limit)
	d.Set("name", r.Name)
	d.Set("playbook", r.Playbook)
	d.Set("project_id", r.Project)
	d.Set("scm_branch", r.ScmBranch)
	d.Set("skip_tags", r.SkipTags)
	d.Set("start_at_task", r.StartAtTask)
	d.Set("survey_enabled", r.SurveyEnabled)
	d.Set("timeout", r.Timeout)
	d.Set("use_fact_cache", r.UseFactCache)
	d.Set("vault_credential_id", r.VaultCredential)
	d.Set("verbosity", r.Verbosity)
	d.Set("webhook_credential_id", r.WebhookCredential)
	d.Set("webhook_service", r.WebhookService)
	extraIDs := getExtraIDs(r)
	d.Set("extra_credential_ids", extraIDs)
	return d
}

func importJobTemplateData(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService

	id, err := strconv.Atoi(d.Id())

	job, err := awxService.GetJobTemplate(id)

	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, err
	}

	resources := []*schema.ResourceData{setJobTemplateResourceData(d, job)}

	return resources, nil
}

func getExtraIDs(template *awxgo.JobTemplate) []int {
	creds := template.SummaryFields.ExtraCredentials
	var ids []int
	for _, c := range creds {
		if cred, ok := c.(*awxgo.Credential); ok {
			ids = append(ids, cred.ID)
		}
	}

	return ids
}
