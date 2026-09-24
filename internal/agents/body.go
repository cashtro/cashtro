package agents

import "strings"

// Organ is one part of the functioning body.
type Organ struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Role  string   `json:"role"`
	Agent string   `json:"agent"`
	Links []string `json:"links"`
}

// Body is the conglomerate. Instinct and Voltron run it.
// The hustler is in it. There is no chair and no dominator.
type Body struct {
	RunBy        []string `json:"runBy"`
	Hustler      bool     `json:"hustler"`
	Chair        string   `json:"chair"`
	Brain        string   `json:"brain"`
	Spy          string   `json:"spy"`
	Enforcer     string   `json:"enforcer"`
	Cyber        string   `json:"cyber"`
	CyberAttacks bool     `json:"cyberAttacks"`
	Organs       []Organ  `json:"organs"`
	Eyes         []string `json:"eyes"`
	Members      []string `json:"members"`
	Connected    bool     `json:"connected"`
}

// FusionBody is the whole conglomerate of agentics, wired as one body.
func FusionBody() Body {
	members := conglomerate()
	eyes := append([]string(nil), members...)
	organs := []Organ{
		{ID: "brain", Name: "Brain", Role: "Instinct decides.", Agent: "instinct", Links: []string{"cia", "fbi", "eyes", "soul"}},
		{ID: "subconscious", Name: "Subconscious", Role: "Memory holds what the brain is not looking at.", Agent: "memory", Links: []string{"brain", "eyes"}},
		{ID: "deep", Name: "Deep conscious", Role: "The Architect thinks underneath. Two ways, then the shorter one.", Agent: "architecte", Links: []string{"brain", "soul"}},
		{ID: "soul", Name: "Soul", Role: "The inside of the brain. Instinct carries it. It is not a second product.", Agent: "instinct", Links: []string{"brain", "deep"}},
		{ID: "eyes", Name: "Eyes", Role: "The eyes see every agentic. They are connected to all of them.", Agent: "eye", Links: eyes},
		{ID: "ears", Name: "Ears", Role: "Explorer listens to the graph and the repos.", Agent: "explorer", Links: []string{"eyes", "brain"}},
		{ID: "mouth", Name: "Mouth", Role: "Comms speaks. Nothing goes out without allow.", Agent: "comms", Links: []string{"brain", "ears"}},
		{ID: "smell", Name: "Smell", Role: "Security smells a secret or a site already live. It does not attack.", Agent: "security", Links: []string{"eyes", "fbi"}},
		{ID: "cia", Name: "CIA", Role: "The Eye is the spy. It gathers and gives that to the brain and to the FBI.", Agent: "eye", Links: []string{"brain", "fbi", "eyes"}},
		{ID: "fbi", Name: "FBI", Role: "The Fusion enforces the laws, punishes with a named story, traces, and secures.", Agent: "fusion", Links: []string{"cia", "brain", "smell"}},
	}
	return Body{
		RunBy:        []string{"instinct", "voltron"},
		Hustler:      true,
		Chair:        "",
		Brain:        "instinct",
		Spy:          "eye",
		Enforcer:     "fusion",
		Cyber:        "security",
		CyberAttacks: false,
		Organs:       organs,
		Eyes:         eyes,
		Members:      members,
		Connected:    covers(eyes, members),
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
	b.WriteString("Instinct is the brain. Voltron keeps the chains on. The hustler is in the conglomerate. Nobody chairs it and nobody dominates it.\n\n")
	b.WriteString("## Roles\n\n")
	b.WriteString("- Brain: instinct. It decides.\n")
	b.WriteString("- CIA: The Eye. It spies and gathers, then gives that information to instinct and to the FBI.\n")
	b.WriteString("- FBI: The Fusion. It enforces every law already written. A failed law stops the gesture, writes the trace, and stores a named sanction story.\n")
	b.WriteString("- Eyes: connected to every agentic. They see the whole conglomerate.\n")
	b.WriteString("- Cybersecurity: security, already in the kernel. It smells secrets and live sites. It does not attack.\n\n")
	b.WriteString("## Body\n\n")
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
