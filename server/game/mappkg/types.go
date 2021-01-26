package mappkg

import (
	"gowog-cloud/Message_proto"
	"gowog-cloud/game/shape"
)

type Map interface {
	ToProto() *Message_proto.Map
	GetWidth() float32
	GetHeight() float32
	GetNumCols() int
	GetNumRows() int
	IsCollide(x float32, y float32) bool
	GetRectBlocks() []shape.Rect
}
