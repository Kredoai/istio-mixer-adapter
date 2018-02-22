package product

import (
	"github.com/apigee/istio-mixer-adapter/apigee/auth"
	"istio.io/istio/mixer/pkg/adapter"
	"strings"
	"net/url"
)

var (
	pm *productManager
)

func Start(baseURL url.URL, log adapter.Logger, env adapter.Env) {
	pm = createProductManager(baseURL, log)
	pm.start(env)
}

func Stop() {
	pm.close()
}

// todo: naive impl, optimize
// todo: check auth scopes
// todo: paths can be wildcards
// 	see: https://docs.apigee.com/developer-services/content/create-api-products#resourcebehavior
func Resolve(ac auth.Context, api, path string) []Details {
	products := pm.getProducts()
	var result []Details
	for _, name := range ac.APIProducts { // find product by name
		apiProduct := products[name]

		for _, attr := range apiProduct.Attributes { // find target services
			if attr.Name == "istio-services" {
				apiProductTargets := strings.Split(attr.Value, ",")
				for _, apiProductTarget := range apiProductTargets { // find target paths
					if apiProductTarget == api {
						validPaths := apiProduct.Resources
						for _, p := range validPaths {
							if p == path {
								result = append(result, apiProduct)
							}
						}
					}
				}
			}
		}
	}
	return result
}

// jwt: application_name -> product list
// jwt: application_name -> authorized scopes
// products -> lookup by id
// 		quota stuff
// 		apiResources (valid paths)
// 		required scopes

// 1. Authenticate & authorize path
// 2. Check quota
//		Get products
//		For each product, check paths
//		For each product w/ a matching path, apply quota?
