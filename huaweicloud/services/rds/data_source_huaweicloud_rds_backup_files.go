// Generated by PMS #553
package rds

import (
	"context"

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

func DataSourceRdsBackupFiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsBackupFilesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the backup ID.`,
			},
			"files": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of backup files.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the file name.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the file size in KB.`,
						},
						"download_link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the link for downloading the backup file.`,
						},
						"link_expired_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the link expiration time.`,
						},
						"database_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the database.`,
						},
					},
				},
			},
			"bucket": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the bucket where the file is located.`,
			},
		},
	}
}

type BackupFilesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newBackupFilesDSWrapper(d *schema.ResourceData, meta interface{}) *BackupFilesDSWrapper {
	return &BackupFilesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRdsBackupFilesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newBackupFilesDSWrapper(d, meta)
	shoBacDowLinRst, err := wrapper.ShowBackupDownloadLink()
	if err != nil {
		return diag.FromErr(err)
	}

	err = wrapper.showBackupDownloadLinkToSchema(shoBacDowLinRst)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	return nil
}

// @API RDS GET /v3/{project_id}/backup-files
func (w *BackupFilesDSWrapper) ShowBackupDownloadLink() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/backup-files"
	params := map[string]any{
		"backup_id": w.Get("backup_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *BackupFilesDSWrapper) showBackupDownloadLinkToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("files", schemas.SliceToList(body.Get("files"),
			func(files gjson.Result) any {
				return map[string]any{
					"name":              files.Get("name").Value(),
					"size":              files.Get("size").Value(),
					"download_link":     files.Get("download_link").Value(),
					"link_expired_time": files.Get("link_expired_time").Value(),
					"database_name":     files.Get("database_name").Value(),
				}
			},
		)),
		d.Set("bucket", body.Get("bucket").Value()),
	)
	return mErr.ErrorOrNil()
}
