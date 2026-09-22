// Package cashtro is the reachable control-plane side of the Graphify tree.
// Fourteen kernel agentics live here. Evolu-Jeunes is the other root.
package cashtro

// Repo is the public manager: github.com/cashtro/cashtro.
func Repo() string { return "cashtro/cashtro" }

// Kernel is Cashtro OS on :8080.
func Kernel() []string {
	return []string{
		Init(), Delivery(), Router(), Research(), Explorer(), Operator(),
		Reviewer(), Architect(), Deploy(), Security(), Memory(), Comms(),
		Planner(), Investigator(),
	}
}

// ControlPlane is recon, the registry, and the Fastify API.
func ControlPlane() string { return "control-plane-api" }

// Init boots the OS and publishes the manifesto. Live.
func Init() string { return "init" }

// Delivery owns idea → concept → production. Live.
func Delivery() string { return "delivery" }

// Router is the optional OpenRouter bus. Resident until a key is bound.
func Router() string { return "router" }

// Research ingests sourced notes. Live.
func Research() string { return "research" }

// Explorer searches processes, ships, and notes. Live.
func Explorer() string { return "explorer" }

// Operator drives browser and desktop. Resident.
func Operator() string { return "operator" }

// Reviewer reads walkthrough artifacts. Resident.
func Reviewer() string { return "reviewer" }

// Architect shapes apps and workflows. Resident.
func Architect() string { return "architect" }

// Deploy owns CI, preview, and promotion. Resident.
func Deploy() string { return "deploy" }

// Security triages findings before a ship. Resident.
func Security() string { return "security" }

// Memory stores and recalls episodic facts. Live.
func Memory() string { return "memory" }

// Comms sends only after a human confirm. Live.
func Comms() string { return "comms" }

// Planner turns a goal into idea-stage ships. Live.
func Planner() string { return "planner" }

// Investigator traces a failing check. Resident.
func Investigator() string { return "investigator" }
