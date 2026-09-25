package agents

import "strings"

// Vein is knowledge in the hustler's veins. The lesson is ours. The book is not copied.
type Vein struct {
	Source string `json:"source"`
	Lesson string `json:"lesson"`
	Copied bool   `json:"copied"`
}

// Organ is one part of the functioning body.
type Organ struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Role  string   `json:"role"`
	Agent string   `json:"agent"`
	Links []string `json:"links"`
}

// EyeGraphic is the map the eyes keep. It grows, and it shows where the body is going.
type EyeGraphic struct {
	Growing bool     `json:"growing"`
	Files   []string `json:"files"`
	Knows   bool     `json:"knows"`
}

// Body is the conglomerate. The hustler chairs it. That chair is the person.
// The hustler relays the ground to instinct. Instinct does the rest.
// Voltron keeps the chains on. This body does not smell and does not eat.
type Body struct {
	RunBy        []string   `json:"runBy"`
	Hustler      bool       `json:"hustler"`
	Person       bool       `json:"person"`
	Chair        string     `json:"chair"`
	Decision     string     `json:"decision"`
	Relay        []string   `json:"relay"`
	RelayTo      string     `json:"relayTo"`
	Vision       bool       `json:"vision"`
	Profile      string     `json:"profile"`
	Ground       bool       `json:"ground"`
	Contacts     bool       `json:"contacts"`
	Veins        []Vein     `json:"veins"`
	Swarm        bool       `json:"swarm"`
	Brain        string     `json:"brain"`
	Spy          string     `json:"spy"`
	Enforcer     string     `json:"enforcer"`
	Cyber        string     `json:"cyber"`
	CyberAttacks bool       `json:"cyberAttacks"`
	Smell        bool       `json:"smell"`
	Eat          bool       `json:"eat"`
	Organs       []Organ    `json:"organs"`
	Eyes         []string   `json:"eyes"`
	Graphic      EyeGraphic `json:"graphic"`
	Members      []string   `json:"members"`
	Connected    bool       `json:"connected"`
}

// FusionBody is the whole conglomerate of agentics, wired as one body.
func FusionBody() Body {
	members := conglomerate()
	eyes := append([]string(nil), members...)
	organs := []Organ{
		{ID: "chair", Name: "Chair", Role: "The hustler is the person. He has the vision, the ground, and the contacts. He relays everything to Einstein.", Agent: "hustler", Links: []string{"instinct", "eyes"}},
		{ID: "brain", Name: "Brain", Role: "Einstein receives the relay from the hustler and does the rest.", Agent: "instinct", Links: []string{"hustler", "cia", "fbi", "eyes", "soul"}},
		{ID: "subconscious", Name: "Subconscious", Role: "Memory holds what the brain is not looking at.", Agent: "memory", Links: []string{"brain", "eyes"}},
		{ID: "deep", Name: "Deep conscious", Role: "The Architect thinks underneath. Two ways, then the shorter one.", Agent: "architecte", Links: []string{"brain", "soul"}},
		{ID: "heart", Name: "Heart", Role: "The Vault is the heart. Core memory, money, and keys. The chain grows. A key slot stores a name, never a secret.", Agent: "memory", Links: []string{"brain", "eyes", "chair"}},
		{ID: "soul", Name: "Soul", Role: "The inside of the brain. Instinct carries it. It is not a second product.", Agent: "instinct", Links: []string{"brain", "deep"}},
		{ID: "eyes", Name: "Eyes", Role: "The eyes see every agentic and draw a graphic that keeps growing, so the body knows where it is going.", Agent: "eye", Links: eyes},
		{ID: "ears", Name: "Ears", Role: "Explorer listens to the graph and the repos.", Agent: "explorer", Links: []string{"eyes", "brain"}},
		{ID: "mouth", Name: "Mouth", Role: "Comms speaks. Nothing goes out without allow.", Agent: "comms", Links: []string{"brain", "ears"}},
		{ID: "cia", Name: "CIA", Role: "The Eye is the spy. It gathers and gives that to the brain and to the FBI.", Agent: "eye", Links: []string{"brain", "fbi", "eyes"}},
		{ID: "fbi", Name: "FBI", Role: "The Fusion enforces the laws. A failed gesture writes the trace, stores the name, saves the story, and retrains the agent until the law, the skill, and the asks are understood.", Agent: "fusion", Links: []string{"cia", "brain", "eyes"}},
	}
	return Body{
		RunBy:        []string{"instinct", "voltron"},
		Hustler:      true,
		Person:       true,
		Chair:        "hustler",
		Decision:     "hustler",
		Relay:        []string{"hustler", "instinct"},
		RelayTo:      "Einstein",
		Vision:       true,
		Profile:      "mogul, 50 Cent",
		Ground:       true,
		Veins:        HustlerVeins(),
		Contacts:     true,
		Swarm:        true,
		Brain:        "instinct",
		Spy:          "eye",
		Enforcer:     "fusion",
		Cyber:        "security",
		CyberAttacks: false,
		Smell:        false,
		Eat:          false,
		Organs:       organs,
		Graphic: EyeGraphic{
			Growing: true,
			Files:   []string{"docs/MAP.md", "inventory/GRAPH.md", "state/blueprint.json"},
			Knows:   true,
		},
		Eyes:      eyes,
		Members:   members,
		Connected: covers(eyes, members),
	}
}

func conglomerate() []string {
	seen := map[string]bool{}
	var out []string
	add := func(id string) {
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		out = append(out, id)
	}
	for _, b := range Brains() {
		add(b.ID)
	}
	for _, m := range Chart().Members() {
		add(m.Agent)
	}
	add("hustler")
	add("security")
	return out
}

// HustlerVeins is the reading that runs through the chair.
// Each lesson is our sentence. None of it is a page from a book.
func HustlerVeins() []Vein {
	return []Vein{
		{Source: "Robert Greene", Lesson: "See the position before the move. Keep the move inside the law.", Copied: false},
		{Source: "The 50th Law", Lesson: "Name the fear, then take the move that fear was blocking.", Copied: false},
		{Source: "50 Cent", Lesson: "The profile is the mogul: vision, ground, contacts, and a long game.", Copied: false},
		{Source: "strategy", Lesson: "One aim, the shorter path, and the reply it opens.", Copied: false},
		{Source: "law", Lesson: "The checklist is already written. A failed gesture retrains the agent.", Copied: false},
		{Source: "business", Lesson: "One offer, one measure, then the next step.", Copied: false},
		{Source: "The Lean Startup", Lesson: "State the bet, measure it, then keep or cut it.", Copied: false},
		{Source: "Good to Great", Lesson: "Name the one thing this body can be best at, and stop the rest.", Copied: false},
		{Source: "Zero to One", Lesson: "Build a step nobody else has, then defend it.", Copied: false},
		{Source: "Competitive Strategy", Lesson: "Know the field, the cost, and the position before the price.", Copied: false},
		{Source: "Profit First", Lesson: "Set the margin aside before the spend.", Copied: false},
		{Source: "college", Lesson: "Study the structure until the agent can say it.", Copied: false},
		{Source: "streets", Lesson: "The ground decides what is real. The chair has lived it.", Copied: false},
		{Source: "math", Lesson: "Count the rate, the expected value, and the tax before the move.", Copied: false},
		{Source: "How to Solve It", Lesson: "Name the unknown, list what is known, then take one step.", Copied: false},
		{Source: "ai-agentic", Lesson: "An agent has a role, a tool, a memory, and a human gate. It does not invent a second product.", Copied: false},
		{Source: "market", Lesson: "Read the market in public. No live order and no token sale without the written AMF answer.", Copied: false},
		{Source: "The Psychology of Money", Lesson: "Time in the position beats a clever trade. The ledger stays internal.", Copied: false},
	}
}

// VaultLesson returns our sentence for one source. The book page is not returned.
func VaultLesson(source string) string {
	for _, vein := range HustlerVeins() {
		if vein.Source == source && !vein.Copied {
			return vein.Lesson
		}
	}
	return ""
}

func covers(eyes, members []string) bool {
	have := map[string]bool{}
	for _, id := range eyes {
		have[id] = true
	}
	if len(eyes) != len(members) {
		return false
	}
	for _, id := range members {
		if !have[id] {
			return false
		}
	}
	return true
}

// FusionResume is the handoff. Rapid fire and code can read it without the kernel.
func FusionResume() string {
	body := FusionBody()
	var b strings.Builder
	b.WriteString("# The Fusion\n\n")
	b.WriteString("Résumé for rapid fire and code. One body. Two GitHubs: Evolu-Jeunes and cashtro.\n\n")
	b.WriteString("## Who runs it\n\n")
	b.WriteString("The hustler is the chair. The hustler is the person. He has the vision. He is a mogul in the profile of 50 Cent: ground, contacts, and the long game. He relays everything to Einstein. Einstein does the rest. Voltron keeps the chains on.\n\n")
	b.WriteString("The person talks to the hustler. The hustler talks back and relays to instinct.\n\n")
	b.WriteString("Every agent and every ID sits on one swarm chain together. The hustler's measure of that swarm is the Miller Research Institute in Silicon Valley, times 3,000 years, at light speed.\n\n")
	b.WriteString("## Roles\n\n")
	b.WriteString("- Chair: hustler. The person. He decides.\n")
	b.WriteString("- Brain: instinct. It receives the relay and does the rest.\n")
	b.WriteString("- CIA: The Eye. It spies and gathers, then gives that information to instinct and to the FBI.\n")
	b.WriteString("- FBI: The Fusion. It enforces every law already written. A failed gesture writes the trace, stores the name, saves the story, and retrains the agent until the law, the skill, and the asks are understood.\n")
	b.WriteString("- Eyes: connected to every agentic. They draw a graphic that keeps growing, from docs/MAP.md, inventory/GRAPH.md, and state/blueprint.json, so the body knows where it is going.\n")
	b.WriteString("- Cybersecurity: security, already in the kernel. It does not attack.\n")
	b.WriteString("- This body does not smell and does not eat.\n\n")
	b.WriteString("## The Vault\n\n")
	b.WriteString("The Vault is the heart. It is the library, chained like a ledger, held in every brain and every cell. Core memory holds the lessons and the book titles. Money holds the internal ledger. Keys hold slot names only. Secret values are not stored. Sections are added slowly.\n\n")
	b.WriteString("## Veins\n\n")
	b.WriteString("This knowledge runs through the hustler. Every strategy and every move has to stand on a vault lesson: business, math, agentic AI, and the market, with the books already named. The lessons are written here. The books are not copied.\n\n")
	for _, vein := range body.Veins {
		b.WriteString("- " + vein.Source + ": " + vein.Lesson + "\n")
	}
	b.WriteString("\n## Body\n\n")
	for _, organ := range body.Organs {
		b.WriteString("- " + organ.Name + " (`" + organ.Agent + "`): " + organ.Role + "\n")
	}
	b.WriteString("\n## Conglomerate\n\n")
	b.WriteString(strings.Join(body.Members, ", ") + "\n\n")
	b.WriteString("## Laws\n\n")
	b.WriteString("The FBI enforces the checklist already built: Loi 25, PIPEDA, CASL, Charte, consumer protection, image, RACJ, AMF. A question becomes a bill. A bill becomes a law only when a rule is passed. The CIA and the FBI each copy the cooperative hierarchy and keep an opposition. The two oppositions scrutinize their own house and build the next structure together.\n\n")
	b.WriteString("## Code\n\n")
	b.WriteString("The control plane is the Go kernel in `cashtro/cashtro`. Call `manager.fusion`. The TypeScript brains stay parallel. There is no import between the two GitHubs.\n")
	return b.String()
}
