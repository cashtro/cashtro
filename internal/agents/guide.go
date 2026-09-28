package agents

// CategoryItem is one member of a desk category.
type CategoryItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Does string `json:"does"`
}

// Category is one shelf on the desk: brain, hustler, and the rest.
type Category struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	What  string         `json:"what"`
	Items []CategoryItem `json:"items"`
}

// Categories is the shelf the desk shows. Each item says what it does.
func Categories() []Category {
	brains := make([]CategoryItem, 0, len(Brains()))
	for _, b := range Brains() {
		brains = append(brains, CategoryItem{ID: b.ID, Name: b.Name, Does: b.Do})
	}
	h := Hustler()
	hs := Hustle()
	lines := make([]CategoryItem, 0, len(Lines()))
	for _, ln := range Lines() {
		lines = append(lines, CategoryItem{ID: ln.ID, Name: ln.Name, Does: ln.Mandate})
	}
	os := TheOS()
	starts := make([]CategoryItem, 0, len(os.Starts))
	for _, start := range os.Starts {
		starts = append(starts, CategoryItem{ID: start.ID, Name: start.Name, Does: start.Task})
	}
	crew := make([]CategoryItem, 0, len(RepairCrew()))
	for _, agent := range RepairCrew() {
		crew = append(crew, CategoryItem{ID: agent.ID, Name: agent.Name, Does: agent.Craft})
	}
	school, _ := LineByID("ecole")
	panda, _ := LineByID("panda")
	home := "Evolu-Jeunes/Panda"
	if len(panda.Repos) > 0 {
		home = panda.Repos[0]
	}
	return []Category{
		{ID: "brain", Name: "Brain", What: "Einstein is the shared brain. Each other brain has one job and closes on him.", Items: brains},
		{ID: "hustler", Name: "Hustler", What: "The hustler is the person. He holds the chair and relays to Einstein.", Items: []CategoryItem{
			{ID: h.ID, Name: "Hustler", Does: "The person. Current: " + h.Current + ". He relays. Einstein does the rest."},
		}},
		{ID: "hustle", Name: "The Hustle", What: "The Hustle is the program. Voltron opens it. Einstein decides.", Items: []CategoryItem{
			{ID: hs.Brain, Name: hs.Name, Does: "Repo " + hs.Repo + ". Opened by " + hs.OpenedBy + ". Decided by " + hs.DecidedBy + "."},
		}},
		{ID: "lines", Name: "Projects", What: "A project is a website or a business. It is not a seat on the agentic chain.", Items: lines},
		{ID: "voltron", Name: "Voltron", What: os.Role + ". A press starts a chain. A website is not a start.", Items: starts},
		{ID: "chain", Name: "Agentic chain", What: "Twenty seats, off the projects. The maître reads the chain. Einstein decides. The plan below is not applied.", Items: crew},
		{ID: "school", Name: "School", What: school.Mandate, Items: []CategoryItem{
			{ID: "panda", Name: panda.Name, Does: "The school that teaches lives in this website, " + home + "."},
			{ID: "educonnexion", Name: "Éduconnexion", Does: "A WordPress website. Not the school. Not a chain."},
		}},
	}
}

// Mark is one reading of the roster. Good means it holds.
type Mark struct {
	ID   string `json:"id"`
	Fact string `json:"fact"`
	Good bool   `json:"good"`
}

// Cut is one chop that waits. Applied stays false until the plan is accepted.
type Cut struct {
	Order   int    `json:"order"`
	Where   string `json:"where"`
	Do      string `json:"do"`
	Applied bool   `json:"applied"`
}

// AgentGuide explains how to build an agent and reads the chain.
// The plan is the list of cuts. The roadmap is later. This desk does not apply either.
type AgentGuide struct {
	How     string     `json:"how"`
	Rules   []string   `json:"rules"`
	Marks   []Mark     `json:"marks"`
	Plan    []Cut      `json:"plan"`
	Roadmap []RoadItem `json:"roadmap"`
	Applied bool       `json:"applied"`
}

// CorrectChain reads the roster and separates what holds from the plan.
func CorrectChain() AgentGuide {
	_, loose := VoltronHolds()
	starts := map[string]bool{}
	for _, start := range TheOS().Starts {
		starts[start.ID] = true
	}
	school, _ := LineByID("ecole")
	pandaOwns := false
	if panda, ok := LineByID("panda"); ok {
		for _, repo := range panda.Repos {
			if repo == "Evolu-Jeunes/Panda" {
				pandaOwns = true
			}
		}
	}
	giantOnTrading := false
	if trading, ok := LineByID("trading"); ok {
		for _, repo := range trading.Repos {
			if repo == "Evolu-Jeunes/Giant" {
				giantOnTrading = true
			}
		}
	}
	fix2Still := false
	if _, ok := Chain("fix2"); ok {
		if _, lineOK := LineByID("fix2"); lineOK {
			fix2Still = true
		}
	}
	crm := CRM()
	fix2Line, _ := LineByID("fix2")
	crmOnFix2 := false
	for _, repo := range fix2Line.Repos {
		if repo == crm.Copy {
			crmOnFix2 = true
		}
	}
	projectSteps := 0
	for _, ln := range Lines() {
		if _, ok := Chain(ln.ID); ok {
			projectSteps++
		}
	}
	marks := []Mark{
		{ID: "edu-website", Fact: "Éduconnexion is a website. It is not a seat on the agentic chain.", Good: !OnChain("educonnexion") && !OnChain("Evolu-Jeunes/educonnexion")},
		{ID: "edu-not-school", Fact: "Éduconnexion stays on the WordPress list and is off the school line.", Good: repoOn("wordpress", "Evolu-Jeunes/educonnexion") && !repoOn("ecole", "Evolu-Jeunes/educonnexion")},
		{ID: "school-in-panda", Fact: "The school teaches inside the Panda website.", Good: pandaOwns && len(school.Repos) == 0},
		{ID: "voltron-starts", Fact: "Voltron starts the agentic chains and does not start a website.", Good: loose == "" && starts["maitre"] && starts["trading"] && !starts["educonnexion"] && !starts["ecole"] && !starts["wordpress"] && !starts["fix2"]},
		{ID: "vapi-panda", Fact: "Vapi is Panda customer service.", Good: Vapi().Project == "panda" && Vapi().Database == "db:panda"},
		{ID: "repos-apart", Fact: "Giant and the trading bot keep two repos.", Good: !giantOnTrading && !GiantSync().Mixed && GiantSync().Together},
		{ID: "panda-pages", Fact: "The extracted Panda backend connects achatRoutes and has no school route. The course pages were not read from the repo.", Good: false},
	}
	scan := ReadPandaScan()
	schoolRoute := false
	for _, section := range scan.Sections {
		if section.ID == "school" && section.Connected {
			schoolRoute = true
		}
	}
	for i := range marks {
		if marks[i].ID == "panda-pages" {
			marks[i].Good = schoolRoute
		}
	}
	plan := []Cut{
		{1, "panda-pages", "The Panda catalog connects purchases in server.js and a demo admin. It does not connect a school route or a price drawer. V1 classes stay on this desk. The repo read waits until the Panda source is here.", false},
		{2, "school", "The school still has its own line, while the teaching lives in the Panda website. Éduconnexion stays a WordPress website. The fold waits until the school route is in the repo.", false},
	}
	if fix2Still {
		plan = append(plan, Cut{len(plan) + 1, "fix2", "Fix Tout keeps its own CRM in its repo. The agent CRM is the other book: knowledge, booking, closing, and a call that stays off. Dashboards land in the Fix Tout repo later.", false})
	}
	if crmOnFix2 {
		plan = append(plan, Cut{len(plan) + 1, "crm", "Evolu-Jeunes/CRM is the copy named on the agent book and also a repo on the Fix Tout line. The agent workflow and the Fix Tout CRM stay two books.", false})
	}
	if projectSteps > 0 {
		plan = append(plan, Cut{len(plan) + 1, "project-steps", "Project lines still carry work steps. They are not seats. The next cut keeps them as project work, or takes them off, when you say which.", false})
	}
	plan = append(plan, Cut{len(plan) + 1, "einstein", "Einstein auto-learns when the analysis goes back to the other agents, through the agent CRM or a direct read. He writes his input after that. The return is on the roadmap and is not called.", false})
	return AgentGuide{
		How: "An agentic is a seat with one job. Name the job, name what it refuses, and name who decides. Einstein decides. Voltron starts a chain that has tasks. A website is not that chain. The fifteen kernel processes stay fifteen at boot. This desk reads the roster, marks what holds, and writes the cuts that do not. It does not apply the plan.",
		Rules: []string{
			"One seat, one job, one refusal.",
			"Einstein is the brain. The hustler relays to him.",
			"Voltron starts an agentic chain. A button is that start.",
			"A website, a project line, and a repo are not seats.",
			"One repo has one owner line.",
			"A send waits on comms.allow. This kernel does not place a live order.",
			"A mark that does not hold stays on the plan until you say to cut.",
		},
		Marks:   marks,
		Plan:    plan,
		Roadmap: Roadmap(),
		Applied: false,
	}
}

func repoOn(line, repo string) bool {
	ln, ok := LineByID(line)
	if !ok {
		return false
	}
	for _, name := range ln.Repos {
		if name == repo {
			return true
		}
	}
	return false
}
