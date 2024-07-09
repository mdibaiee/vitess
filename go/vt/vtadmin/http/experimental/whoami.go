package experimental

import (
	"context"

	vtadminhttp "mdibaiee/vitess/oracle/go/vt/vtadmin/http"
	"mdibaiee/vitess/oracle/go/vt/vtadmin/rbac"
)

// WhoAmI is an experimental route for extracting authenticated Actors from
// the request, to see who is authenticated on the frontend.
func WhoAmI(ctx context.Context, r vtadminhttp.Request, api *vtadminhttp.API) *vtadminhttp.JSONResponse {
	data := map[string]any{}
	actor, ok := rbac.FromContext(ctx)
	data["authenticated"] = ok
	if ok {
		data["actor"] = actor
	}

	return vtadminhttp.NewJSONResponse(data, nil)
}
