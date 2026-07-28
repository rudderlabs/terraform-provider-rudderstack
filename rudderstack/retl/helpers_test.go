package retl_test

// ptr returns a pointer to v. Used for rudder-iac request fields that became
// pointers in v0.19.0+ (e.g. CreateRETLConnectionRequest.SyncBehaviour).
func ptr[T any](v T) *T { return &v }
