// Generated by PMS #682
package aad

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
)

func DataSourceAadDomainCertificate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAadDomainCertificateRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the AAD domain ID.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name.`,
			},
			"cert_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The certificate information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The certificate expiration time.`,
						},
						"expire_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The certificate expiration status.`,
						},
						"cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate name.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The certificate ID.`,
						},
						"apply_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name that the certificate applies to.`,
						},
					},
				},
			},
		},
	}
}

type DomainCertificateDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newDomainCertificateDSWrapper(d *schema.ResourceData, meta interface{}) *DomainCertificateDSWrapper {
	return &DomainCertificateDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceAadDomainCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newDomainCertificateDSWrapper(d, meta)
	shoDomCerRst, err := wrapper.ShowDomainCertificate()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showDomainCertificateToSchema(shoDomCerRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API AAD GET /v2/aad/domains/{domain_id}/certificate
func (w *DomainCertificateDSWrapper) ShowDomainCertificate() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "aad")
	if err != nil {
		return nil, err
	}

	uri := "/v2/aad/domains/{domain_id}/certificate"
	uri = strings.ReplaceAll(uri, "{domain_id}", w.Get("domain_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *DomainCertificateDSWrapper) showDomainCertificateToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("domain_name", body.Get("domain_name").Value()),
		d.Set("cert_info", schemas.ObjectToList(body.Get("cert_info"),
			func(certInfo gjson.Result) any {
				return map[string]any{
					"expire_time":   certInfo.Get("expire_time").Value(),
					"expire_status": certInfo.Get("expire_status").Value(),
					"cert_name":     certInfo.Get("cert_name").Value(),
					"id":            certInfo.Get("id").Value(),
					"apply_domain":  certInfo.Get("apply_domain").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
