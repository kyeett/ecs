package system

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
)

// Follow is responsible for moving an entity following another
type Follow struct {
	em  *entity.Manager
	log logging.Logger
}

// NewFollow creates a new Follow system
func NewFollow(em *entity.Manager, logger logging.Logger) *Follow {
	return &Follow{
		em:  em,
		log: logger.WithField("s", "follow"),
	}
}

// Update the Follow system
func (i *Follow) Update(dt float64) {
	for _, e := range i.em.FilteredEntities(components.PosType, components.FollowingType) {
		follow := i.em.Follow(e)
		pos := i.em.Pos(e)
		posFollowed := i.em.Pos(follow.ID)
		pos.Vec = posFollowed.Vec
		i.log.Debugf("%q following %q", e, follow.ID)
	}
}
