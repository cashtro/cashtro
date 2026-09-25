package agents

// TrinityChapter is one step of the agency shape.
// Existing agents come first. Departments come after. One brain. One runner.
type TrinityChapter struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Here  string `json:"here"`
}

// Trinity is this kernel as one instance. It is not a second product.
type Trinity struct {
	Name         string           `json:"name"`
	Instance     string           `json:"instance"`
	SharedBrain  string           `json:"sharedBrain"`
	Orchestrator string           `json:"orchestrator"`
	Workspace    string           `json:"workspace"`
	Commercial   string           `json:"commercial"`
	Chapters     []TrinityChapter `json:"chapters"`
}

// TheTrinity is the six-part shape. Agents already registered, then the departments.
func TheTrinity() Trinity {
	return Trinity{
		Name:         "Trinity",
		Instance:     "cashtro/cashtro",
		SharedBrain:  "instinct",
		Orchestrator: "voltron",
		Workspace:    "http://localhost:8080",
		Commercial:   "ops",
		Chapters: []TrinityChapter{
			{ID: "what", Title: "What Trinity is, and your own instance running", Here: "This kernel. One local instance. go run ./cmd/cashtro."},
			{ID: "agents", Title: "Bring your existing agents, then build the departments", Here: "The registered agentics come first. Chart departments come after."},
			{ID: "brain", Title: "Give your agency one shared brain", Here: "Einstein is the one shared brain. The hustler relays to him."},
			{ID: "orchestration", Title: "Orchestration: one agent to run them all", Here: "Voltron keeps every chain on."},
			{ID: "workspace", Title: "Workspace: the interface for your team", Here: "The desk at http://localhost:8080."},
			{ID: "operations", Title: "Operations and the commercial model", Here: "OPS runs the work. Each division has a product, a revenue, and a guard."},
		},
	}
}
