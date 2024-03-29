package poi

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"hash/fnv"

	"google.golang.org/protobuf/proto"
)

var poi = []string{
	"天安门",
	"广州塔",
	"庐江",
	"体育西",
	"广州南",
	"广州北",
}

type Manager struct {
}

func (*Manager) Resolve(c context.Context, loc *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(loc)
	if err != nil {
		return "", err
	}

	h := fnv.New32()
	h.Write(b)

	return poi[int(h.Sum32())%len(poi)], nil
}
