package domain

type User struct {
	ID      int64
	Code    string
	Budget  *float64
	Updated int64
}
