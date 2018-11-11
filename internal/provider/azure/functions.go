package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/twistlock/cloud-discovery/internal/shared"
	"strings"
)

// DiscoverFunctions retrieves all Azure functions
func DiscoverFunctions(opt Options, emitFn func(result shared.CloudDiscoveryResult)) error {
	spt, err := spt(opt)
	if err != nil {
		return err
	}
	ac := web.NewAppsClient(opt.SubscriptionID)
	ac.Authorizer = autorest.NewBearerAuthorizer(spt)

	functionApps, err := ac.List(context.Background())
	if err != nil {
		return err
	}
	var result shared.CloudDiscoveryResult
	result.Type = "Azure Functions"
	const azureFunctionKind = "functionapp"
	for {
		apps := functionApps.Values()
		if len(apps) == 0 {
			return err
		}
		for _, app := range apps {
			if !strings.HasPrefix(strings.ToLower(*app.Kind), azureFunctionKind) {
				continue
			}

			if app.Name == nil || app.Location == nil {
				continue
			}

			emitFn(shared.CloudDiscoveryResult{
				Region: *app.Location,
				Type:   "Azure Functions",
				Assets: []shared.CloudAsset{
					{ID: *app.Name, Data: app},
				},
			})
		}

		if err := functionApps.Next(); err != nil {
			return err
		}
	}
}

// spt returns an authenticated service principal token using provided credentials
func spt(opt Options) (*adal.ServicePrincipalToken, error) {
	oauthConfig, err := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, opt.TenantID)
	if err != nil {
		return nil, err
	}
	spt, err := adal.NewServicePrincipalToken(*oauthConfig, opt.ClientID, opt.Secret, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		return nil, err
	}
	return spt, nil
}
