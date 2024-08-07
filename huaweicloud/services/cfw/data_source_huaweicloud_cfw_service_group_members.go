// Generated by PMS #159
package cfw

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCfwServiceGroupMembers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCfwServiceGroupMembersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the service group ID.`,
			},
			"key_word": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the key word.`,
			},
			"group_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the service group type.`,
			},
			"item_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the service group member ID.`,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the firewall instance ID.`,
			},
			"protocol": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the protocol type.`,
			},
			"source_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source port.`,
			},
			"dest_port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the destination port.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The service group member list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"item_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The service group member ID.`,
						},
						"protocol": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The protocol type.`,
						},
						"source_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The source port.`,
						},
						"dest_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The destination port.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The service group member description.`,
						},
					},
				},
			},
		},
	}
}

type ServiceGroupMembersDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newServiceGroupMembersDSWrapper(d *schema.ResourceData, meta interface{}) *ServiceGroupMembersDSWrapper {
	return &ServiceGroupMembersDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCfwServiceGroupMembersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newServiceGroupMembersDSWrapper(d, meta)
	listServiceItemsRst, err := wrapper.ListServiceItems()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listServiceItemsToSchema(listServiceItemsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CFW GET /v1/{project_id}/service-items
func (w *ServiceGroupMembersDSWrapper) ListServiceItems() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cfw")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/service-items"
	params := map[string]any{
		"set_id":                 w.Get("group_id"),
		"key_word":               w.Get("key_word"),
		"fw_instance_id":         w.Get("fw_instance_id"),
		"query_service_set_type": w.Get("group_type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("data.records", "offset", "limit", 1024).
		Filter(
			filters.New().From("data.records").
				Where("source_port", "=", w.Get("source_port")).
				Where("dest_port", "=", w.Get("dest_port")).
				Where("protocol", "=", w.Get("protocol")).
				Where("item_id", "=", w.Get("item_id")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *ServiceGroupMembersDSWrapper) listServiceItemsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("records", schemas.SliceToList(body.Get("data.records"),
			func(records gjson.Result) any {
				return map[string]any{
					"item_id":     records.Get("item_id").Value(),
					"protocol":    records.Get("protocol").Value(),
					"source_port": records.Get("source_port").Value(),
					"dest_port":   records.Get("dest_port").Value(),
					"description": records.Get("description").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
