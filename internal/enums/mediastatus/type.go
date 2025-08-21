package mediastatus

type Type string

const (
	Pending   Type = "pending"
	Validated Type = "validated"
	Rejected  Type = "rejected"
)
