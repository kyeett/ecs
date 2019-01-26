package system

import (
	"github.com/kyeett/ecs/entity"
	"github.com/kyeett/ecs/events"
	"github.com/kyeett/ecs/logging"
	"github.com/kyeett/gomponents/components"
	"github.com/peterhellberg/gfx"
)

// Parenting
type Parenting struct {
	em    *entity.Manager
	outCh chan events.Event
	log   logging.Logger
}

// NewParenting creates a new Parenting system
func NewParenting(em *entity.Manager, ch chan events.Event, logger logging.Logger) *Parenting {
	return &Parenting{
		em:    em,
		outCh: ch,
		log:   logger.WithField("s", "parenting"),
	}
}

// Update the Parenting system
func (r *Parenting) Update(dt float64) {

	// Check if parented entities has moved away from their parent horizontally
	for _, e := range r.em.FilteredEntities(components.ParentedType, components.HitboxType, components.PosType) {
		parented := r.em.Parented(e)
		hb := movedHitbox(e, r.em)
		hbParent := resolvRectangle(movedHitbox(parented.ID, r.em))

		bottomRight := hb.Max

		// Create 3 rays, left, center, right of bottom. Move into entity hitbox above to make more robust
		rays := []gfx.Vec{
			bottomRight.Sub(gfx.V(hb.W(), -5)),
			bottomRight.Sub(gfx.V(hb.W()/2, -5)),
			bottomRight.Sub(gfx.V(0, -5))}

		isColliding := false
		for _, v := range rays {
			l := resolvLine(v, v.Add(gfx.V(0, 10)))
			if l.IsColliding(hbParent) {
				isColliding = true
				break
			}
		}

		// Rays don't intersect parent, break parenting
		if isColliding == false {
			r.em.Remove(e, components.ParentedType)
		}
	}
}

func movedHitbox(e string, em *entity.Manager) gfx.Rect {
	pos := em.Pos(e)
	hb := em.Hitbox(e)
	return hb.Rect.Moved(pos.Vec)
}

// Send is an empty method to implement the System interface
func (r *Parenting) Send(ev events.Event) {}