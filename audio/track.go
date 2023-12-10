package audio

import "github.com/pion/interceptor"

type Track interface {
	Read([]byte) (int, interceptor.Attributes, error)
}
