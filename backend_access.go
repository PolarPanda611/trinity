package trinity

//DefaultAccessBackend for Mixin
func DefaultAccessBackend(v *ViewSetRunTime) error {
	//compare permission list
	userPermission := v.Gcontext.GetStringSlice("UserPermission")
	accessBackendRequire := v.AccessBackendRequire
	return CheckAccessAuthorization(accessBackendRequire, userPermission)
	// return nil, nil
}

// CheckAccessAuthorization to check access authorization
func CheckAccessAuthorization(requiredPermission, userPermission []string) error {
	if SliceInSlice(requiredPermission, userPermission) {
		return nil
	}
	return ErrAccessAuthCheckFailed
}
