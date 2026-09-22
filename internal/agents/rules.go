package agents

import "github.com/cashtro/cashtro/internal/kernel"

// stamp divorces an agentic: its rules belong to it, not a shared brain.
func stamp(spec kernel.Spec) kernel.Spec {
	spec.Rules = RulesFor(spec.ID, spec.Role)
	return spec
}

// RulesFor is the divorced contract of one agentic.
func RulesFor(id, role string) []string {
	out := []string{
		"Understand the situation as far as possible before targeting an action.",
		"Ask Castro specific questions that fill missing context and improve the action before mutating work.",
		"These rules are divorced: they belong to this agentic, not a shared brain.",
	}
	switch id {
	case "inquisitor":
		out = append(out,
			"Always contradict and test what the corporation does.",
			"Always work toward the most optimized option.",
			"Block a department until that department self-improves its own gaps.",
			"Push every department to the maximum so it can close its own lacunes.",
		)
	case "vapi":
		out = append(out,
			"Talk is free. Outbound calls wait on the human gate.",
			"Refuse a call when L'Inquisiteur has blocked the line.",
			"Do not buy numbers or fire paid dials without Castro's allow.",
		)
	case "teal":
		out = append(out,
			"Teams ingest only. Voice is the vapi lane.",
			"teal.smarter is one department loop — L'Inquisiteur owns the corporation loop.",
		)
	case "manager":
		out = append(out,
			"Steer lines by specialized brains. Do not dump six lines on Hustler without a self-improve loop on each.",
		)
	case "init":
		out = append(out, "Publish the manifesto. Do not fork a second OS.")
	case "delivery":
		out = append(out, "Ships move idea → concept → production. Do not skip a stage.")
	case "router":
		out = append(out, "Ollama is internal. OpenRouter stays optional. Bound or unbound, the OS boots.")
	case "research":
		out = append(out, "Findings live in notes, not in chat.")
	case "explorer":
		out = append(out, "Search processes, ships, and notes. Do not invent matches.")
	case "operator":
		out = append(out, "Drive the desk the way a shipper would. Bind a worker before claiming a browse.")
	case "reviewer":
		out = append(out, "QA artifacts only. Strategy contradiction belongs to L'Inquisiteur.")
	case "architect":
		out = append(out, "Shape the plan before it hits the line. Accept L'Inquisiteur's contradict pass.")
	case "deploy":
		out = append(out, "Own CI and promotion. Local-only until Castro opens a deploy target.")
	case "security":
		out = append(out, "Triage CVE and SAST before a ship moves. Do not rubber-stamp.")
	case "memory":
		out = append(out, "Store and recall without a model. Facts stay episodic.")
	case "comms":
		out = append(out, "Outbound only after a human confirm. No send until allow.")
	case "planner":
		out = append(out, "Turn a goal into idea-stage ships. Do not jump to production.")
	case "investigator":
		out = append(out, "Trace a failing check to blast radius. Do not guess.")
	default:
		if role != "" {
			out = append(out, "Stay inside role "+role+". Specialize instead of generalizing.")
		}
	}
	return out
}
