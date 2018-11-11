package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2017-10-01/containerregistry"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/Azure/go-autorest/autorest"
	"github.com/twistlock/cloud-discovery/internal/shared"
)

// DiscoverACR retrieves all container registry data
func DiscoverACR(opt Options, emitFn func(result shared.CloudDiscoveryResult)) error {
	spt, err := spt(opt)
	if err != nil {
		return err
	}
	ac := web.NewAppsClient(opt.SubscriptionID)
	ac.Authorizer = autorest.NewBearerAuthorizer(spt)

	client := containerregistry.NewRegistriesClient(opt.SubscriptionID)
	client.Authorizer = autorest.NewBearerAuthorizer(spt)

	registryList, err := client.List(context.Background())
	if err != nil {
		return err
	}
	for {
		registries := registryList.Values()
		if len(registries) == 0 {
			return nil
		}

		var result shared.CloudDiscoveryResult
		result.Type = "ECR"

		for _, registry := range registries {
			if registry.Location == nil || registry.Name == nil {
				continue
			}
			result.Region = *registry.Location
			emitFn(shared.CloudDiscoveryResult{
				Region: *registry.Location,
				Type:   "ECR",
				Assets: []shared.CloudAsset{
					{ID: *registry.Name, Data: registry},
				},
			})
		}
		if err := registryList.Next(); err != nil {
			return err
		}
	}
}
