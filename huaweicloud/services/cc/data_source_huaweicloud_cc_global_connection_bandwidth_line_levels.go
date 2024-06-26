// Generated by PMS #110
package cc

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

func DataSourceCcGlobalConnectionBandwidthLineLevels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcGlobalConnectionBandwidthLineLevelsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"line_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Line ID.`,
			},
			"local_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Local access point code included in the line specification.`,
			},
			"remote_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Remote access point code included in the line specification.`,
			},
			"levels": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Line grade.`,
			},
			"line_levels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The line grade list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The line ID.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the line was created.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Time when the line was updated.`,
						},
						"local_area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Local access point.`,
						},
						"remote_area": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Remote access point.`,
						},
						"levels": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Line grade.`,
						},
					},
				},
			},
		},
	}
}

type GlobalConnectionBandwidthLineLevelsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newGlobalConnectionBandwidthLineLevelsDSWrapper(d *schema.ResourceData, meta interface{}) *GlobalConnectionBandwidthLineLevelsDSWrapper {
	return &GlobalConnectionBandwidthLineLevelsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCcGlobalConnectionBandwidthLineLevelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newGlobalConnectionBandwidthLineLevelsDSWrapper(d, meta)
	lgcbllr, err := wrapper.ListGlobalConnectionBandwidthLineLevels()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.listGlobalConnectionBandwidthLineLevelsToSchema(lgcbllr)
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	return nil
}

// @API CC GET /v3/{domain_id}/gcb/line-levels
func (w *GlobalConnectionBandwidthLineLevelsDSWrapper) ListGlobalConnectionBandwidthLineLevels() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cc")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{domain_id}/gcb/line-levels"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Config.DomainID)
	params := map[string]any{
		"id":          w.Get("line_id"),
		"local_area":  w.Get("local_area"),
		"remote_area": w.Get("remote_area"),
		"levels":      w.Get("levels"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("line_levels", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *GlobalConnectionBandwidthLineLevelsDSWrapper) listGlobalConnectionBandwidthLineLevelsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("line_levels", schemas.SliceToList(body.Get("line_levels"),
			func(linLev gjson.Result) any {
				return map[string]any{
					"id":          linLev.Get("id").Value(),
					"created_at":  linLev.Get("created_at").Value(),
					"updated_at":  linLev.Get("updated_at").Value(),
					"local_area":  linLev.Get("local_area").Value(),
					"remote_area": linLev.Get("remote_area").Value(),
					"levels":      schemas.SliceToStrList(linLev.Get("levels")),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
