// Package castro is the root of the Graphify tree: one person, two homes.
package castro

import (
	"github.com/cashtro/cashtro/internal/graphifytree/cashtro"
	"github.com/cashtro/cashtro/internal/graphifytree/evolujeunes"
)

// Home is one GitHub side Castro owns.
type Home struct {
	Name  string
	Nodes []string
}

// Map is the whole tree: Cashtro agentics and the Evolu-Jeunes fleet.
func Map() (cashtroHome, evoluHome Home) {
	return OwnsCashtro(), OwnsEvoluJeunes()
}

// OwnsCashtro is the reachable control plane and its fourteen agentics.
func OwnsCashtro() Home {
	return Home{
		Name:  cashtro.Repo(),
		Nodes: append([]string{cashtro.ControlPlane()}, cashtro.Kernel()...),
	}
}

// OwnsEvoluJeunes is the delivery org, the live site, and catalog ships.
func OwnsEvoluJeunes() Home {
	nodes := []string{evolujeunes.Org(), evolujeunes.Member(), evolujeunes.Site(), evolujeunes.PrivateRepos()}
	nodes = append(nodes, evolujeunes.Fleet()...)
	// Delivery, explorer, and planner are the kernel seats that touch ships.
	_ = cashtro.Delivery()
	_ = cashtro.Explorer()
	_ = cashtro.Planner()
	return Home{Name: evolujeunes.Org(), Nodes: nodes}
}
