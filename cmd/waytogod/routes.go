// +build !experimental

package main

import "github.com/kmova/waytogo/api/server/router"

func addExperimentalRouters(routers []router.Router) []router.Router {
	return routers
}
