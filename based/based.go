package based

type grp struct {
	ID int64
	name string
}
type usr struct {
	ID int64
	name string
}
type token struct {
	ID int64
	user_ID int64
	expiry int64
	max int16
	used int16
}
type permissions struct {
	ID int64
	resource_path string
	allowed bool
	apply_recursive bool
}
type requests struct {
	ID int64
	ip string
	access_time int64
	resource_path string
	token string
}