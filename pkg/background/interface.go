package background

import "context"

type Background interface {
	Start(context.Context) error
}
