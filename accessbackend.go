package trinity

//DefaultAccessBackend for Mixin
func DefaultAccessBackend(v *ViewSetRunTime) error {
	//compare permission list
	userPermission := v.Gcontext.GetStringSlice("UserPermission")
	accessBackendRequire := v.AccessBackendRequire
	return CheckAccessAuthorization(accessBackendRequire, userPermission)
	// return nil
}
