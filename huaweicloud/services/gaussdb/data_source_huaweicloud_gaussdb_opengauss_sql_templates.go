// Generated by PMS #541
package gaussdb

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

func DataSourceGaussdbOpengaussSqlTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussdbOpengaussSqlTemplatesRead,

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
				Description: `Specifies the ID of a GaussDB OpenGauss instance.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a GaussDB OpenGauss instance node.`,
			},
			"sql_model": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the SQL template.`,
			},
			"node_limit_sql_model_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the information about the SQL template for SQL throttling.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sql_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the SQL ID of the throttling task.`,
						},
						"sql_model": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the SQL template of the throttling task.`,
						},
					},
				},
			},
		},
	}
}

type OpengaussSqlTemplatesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newOpengaussSqlTemplatesDSWrapper(d *schema.ResourceData, meta interface{}) *OpengaussSqlTemplatesDSWrapper {
	return &OpengaussSqlTemplatesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceGaussdbOpengaussSqlTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newOpengaussSqlTemplatesDSWrapper(d, meta)
	lisNodLimSqlModRst, err := wrapper.ListNodeLimitSqlModel()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listNodeLimitSqlModelToSchema(lisNodLimSqlModRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/list-node-limit-sql-model
func (w *OpengaussSqlTemplatesDSWrapper) ListNodeLimitSqlModel() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "opengauss")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/instances/{instance_id}/list-node-limit-sql-model"
	uri = strings.ReplaceAll(uri, "{instance_id}", w.Get("instance_id").(string))
	params := map[string]any{
		"node_id":   w.Get("node_id"),
		"sql_model": w.Get("sql_model"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("node_limit_sql_model_list", "offset", "limit", 0).
		Request().
		Result()
}

func (w *OpengaussSqlTemplatesDSWrapper) listNodeLimitSqlModelToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("node_limit_sql_model_list", schemas.SliceToList(body.Get("node_limit_sql_model_list"),
			func(nlsml gjson.Result) any {
				return map[string]any{
					"sql_id":    nlsml.Get("sql_id").Value(),
					"sql_model": nlsml.Get("sql_model").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
