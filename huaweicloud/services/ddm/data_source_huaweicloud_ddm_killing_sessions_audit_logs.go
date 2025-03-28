// Generated by PMS #446
package ddm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDdmKillingSessionsAuditLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdmKillingSessionsAuditLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DDM instance ID or ID of the associated RDS instance.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time in UTC, accurate to milliseconds.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time in UTC, accurate to milliseconds.`,
			},
			"process_audit_logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of killing sessions audit log.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"process_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the session ID.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the DDM instance ID.`,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the DDM instance name.`,
						},
						"execute_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the execute time in UTC format.`,
						},
						"execute_user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the user who executes the task.`,
						},
					},
				},
			},
		},
	}
}

type KillingSessionsAuditLogsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newKillingSessionsAuditLogsDSWrapper(d *schema.ResourceData, meta interface{}) *KillingSessionsAuditLogsDSWrapper {
	return &KillingSessionsAuditLogsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDdmKillingSessionsAuditLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newKillingSessionsAuditLogsDSWrapper(d, meta)
	shoProAudLogRst, err := wrapper.ShowProcessesAuditLog()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showProcessesAuditLogToSchema(shoProAudLogRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DDM GET /v3/{project_id}/instances/{instance_id}/processes-audit-log
func (w *KillingSessionsAuditLogsDSWrapper) ShowProcessesAuditLog() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ddm")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/instances/{instance_id}/processes-audit-log"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"start_time": w.Get("start_time"),
		"end_time":   w.Get("end_time"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("process_audit_logs", "offset", "limit", 100).
		Request().
		Result()
}

func (w *KillingSessionsAuditLogsDSWrapper) showProcessesAuditLogToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("process_audit_logs", schemas.SliceToList(body.Get("process_audit_logs"),
			func(proAudLogs gjson.Result) any {
				return map[string]any{
					"process_id":        proAudLogs.Get("process_id").Value(),
					"instance_id":       proAudLogs.Get("instance_id").Value(),
					"instance_name":     proAudLogs.Get("instance_name").Value(),
					"execute_time":      proAudLogs.Get("execute_time").Value(),
					"execute_user_name": proAudLogs.Get("execute_user_name").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
