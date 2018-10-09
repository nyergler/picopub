//go:generate jsonenums -type=PostType
//go:generate stringer -type=PostType

package picopub

type PostType int

const (
	Unknown PostType = iota
	Entry
	Read
	Like
)

type Create struct {
	Type       string
	Properties map[string]string
}

type Update struct{}
