package agents

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Board is the one operating picture. The kernel boots from this,
// not from hardcoded counts.
type Board struct {
	Version    string        `json:"version"`
	Fiches     int           `json:"fiches"`
	Repos      []RepoCard    `json:"repos"`
	Lines      []LineView    `json:"lines"`
	Links      []Link        `json:"links"`
	Agents     []string      `json:"agents"`
	Org        Organization  `json:"org"`
	Brains     []LoadedBrain `json:"brains"`
	Compliance Compliance    `json:"compliance"`
	Graph      GraphView     `json:"graph"`
	Gaps       []string      `json:"gaps"`
	Rules      []string      `json:"rules"`
}

// RepoCard is one GitHub repo as the fiches describe it.
type RepoCard struct {
	FullName string `json:"fullName"`
	Org      string `json:"org"`
	Language string `json:"language,omitempty"`
	Brain    string `json:"brain"`
	Line     string `json:"line"`
	LineID   string `json:"lineId,omitempty"`
	Status   string `json:"status"`
}

// LineView is a business line with the repos actually filed under it.
type LineView struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Brain string   `json:"brain"`
	Repos []string `json:"repos"`
}

// GraphView is the accounts map, not the 40k-node extract.
type GraphView struct {
	RepoNodes int `json:"repoNodes"`
	Bridges   int `json:"bridges"`
}

var (
	reRepo = regexp.MustCompile(`\| Repo \| ` + "`" + `([^` + "`" + `]+)` + "`" + ` \|`)
	reCell = regexp.MustCompile(`\| (Langage|Cerveau) \| ([^|]+) \|`)
	reLine = regexp.MustCompile(`\*\*Ligne :\*\* (.+)`)
	reStat = regexp.MustCompile(`\[x\] (idee|concept|production)`)
)

// LoadBoard reads state/operating.json by walking up to the module root.
func LoadBoard() (Board, error) {
	path, err := findFile("state", "operating.json")
	if err != nil {
		return boardFromLines(), nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return boardFromLines(), nil
	}
	var b Board
	if err := json.Unmarshal(raw, &b); err != nil {
		return Board{}, err
	}
	return AlignBoard(b), nil
}

// AlignBoard keeps the fiche cards and puts lines, links, and the chart
// back on the code roster. WordPress themes leave the other department.
func AlignBoard(b Board) Board {
	b.Version = "2"
	b.Org = Chart()
	b.Brains = Brains()
	b.Compliance = Loi()
	b.Agents = agentNames(b.Org)
	b.Links = Links()
	b.Repos = retargetCards(b.Repos)
	b.Lines = tally(b.Repos)
	fresh := boardFromLines().Rules
	if len(b.Rules) == len(fresh) && len(b.Rules) > 0 {
		b.Rules[0] = fresh[0]
	} else {
		b.Rules = fresh
	}
	return b
}

// BuildBoard syncs fiches, the line roster, and the accounts graph into one board.
func BuildBoard(ficheDir, graphPath string) (Board, error) {
	cards, err := readFiches(ficheDir)
	if err != nil {
		return Board{}, err
	}
	b := boardFromLines()
	b.Repos = retargetCards(cards)
	b.Fiches = len(cards)
	b.Lines = tally(b.Repos)
	b.Gaps = gaps(cards)
	if graphPath != "" {
		g, err := readGraph(graphPath)
		if err != nil {
			return Board{}, err
		}
		b.Graph = g
	}
	return b, nil
}

// SaveBoard writes the board as indented JSON.
func SaveBoard(path string, b Board) error {
	raw, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}

func boardFromLines() Board {
	views := make([]LineView, 0, len(Lines()))
	for _, ln := range Lines() {
		views = append(views, LineView{ID: ln.ID, Name: ln.Name, Brain: ln.Brain, Repos: append([]string(nil), ln.Repos...)})
	}
	org := Chart()
	return Board{
		Version:    "2",
		Lines:      views,
		Links:      Links(),
		Agents:     agentNames(org),
		Org:        org,
		Brains:     Brains(),
		Compliance: Loi(),
		Rules: []string{
			"Coopérative : CEO, CTO, CMP. Neuf départements. Quinze employés spécialisés. Le flow stack relie chaque division à une fonction et à un produit.",
			"Avant chaque coup : une question, la position, le coup de pouvoir, la réponse adverse.",
			"La chaîne des coups ne se réécrit pas.",
			"L'agence centrale rapporte. Epicenter Einstein décide. Elle n'est pas un siège de ce tableau.",
			"Le roster de 1001 n'est pas des processus.",
			"Aucun déploiement, aucun envoi, aucun ordre de trading sans comms.allow.",
			"Les sites clients déjà en production se lisent, ils ne se déploient pas.",
			"Chaque division a un revenu et une garde. On ne pousse pas le volume sans mesure.",
			"Un build n'est pas fini tant que les portes de manager.loi ne sont pas passées.",
		},
	}
}

func agentNames(org Organization) []string {
	members := org.Members()
	out := make([]string, 0, len(members)+1)
	for _, m := range members {
		out = append(out, m.Agent)
	}
	out = append(out, "fusion")
	return out
}

func readFiches(dir string) ([]RepoCard, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	byName := map[string]string{}
	for _, ln := range Lines() {
		byName[ln.Name] = ln.ID
	}
	var cards []RepoCard
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		text := string(raw)
		full := ""
		if m := reRepo.FindStringSubmatch(text); len(m) == 2 {
			full = m[1]
		}
		if full == "" {
			full = strings.Replace(strings.TrimSuffix(e.Name(), ".md"), "__", "/", 1)
		}
		org := full
		if i := strings.Index(full, "/"); i >= 0 {
			org = full[:i]
		}
		brain, lang := "", ""
		for _, m := range reCell.FindAllStringSubmatch(text, -1) {
			val := strings.TrimSpace(m[2])
			if m[1] == "Cerveau" {
				brain = val
			} else {
				lang = val
			}
		}
		lineName := ""
		if m := reLine.FindStringSubmatch(text); len(m) == 2 {
			lineName = strings.TrimSpace(m[1])
		}
		status := "inconnu"
		if m := reStat.FindStringSubmatch(text); len(m) == 2 {
			status = m[1]
		}
		cards = append(cards, RepoCard{
			FullName: full,
			Org:      org,
			Language: lang,
			Brain:    brain,
			Line:     lineName,
			LineID:   byName[lineName],
			Status:   status,
		})
	}
	sort.Slice(cards, func(i, j int) bool { return cards[i].FullName < cards[j].FullName })
	return cards, nil
}

// rosterOwner is the department that holds this repo.
// A WordPress theme stays on that line even when an old fiche still says Proximity.
func rosterOwner(full string) (Line, bool) {
	var found Line
	ok := false
	for _, ln := range Lines() {
		for _, repo := range ln.Repos {
			if repo != full {
				continue
			}
			if ln.ID == "wordpress" {
				return ln, true
			}
			if !ok {
				found = ln
				ok = true
			}
		}
	}
	return found, ok
}

func retargetCards(cards []RepoCard) []RepoCard {
	out := append([]RepoCard(nil), cards...)
	for i, c := range out {
		ln, ok := rosterOwner(c.FullName)
		if !ok {
			continue
		}
		out[i].Line = ln.Name
		out[i].LineID = ln.ID
	}
	return out
}

func tally(cards []RepoCard) []LineView {
	present := map[string]bool{}
	for _, c := range cards {
		present[c.FullName] = true
	}
	useAll := len(cards) == 0
	onRoster := map[string]bool{}
	buckets := map[string][]string{}
	for _, ln := range Lines() {
		for _, repo := range ln.Repos {
			onRoster[repo] = true
			if useAll || present[repo] {
				buckets[ln.ID] = append(buckets[ln.ID], repo)
			}
		}
	}
	for _, c := range cards {
		if onRoster[c.FullName] || c.LineID == "" {
			continue
		}
		buckets[c.LineID] = append(buckets[c.LineID], c.FullName)
	}
	views := make([]LineView, 0, len(Lines()))
	for _, ln := range Lines() {
		repos := buckets[ln.ID]
		sort.Strings(repos)
		views = append(views, LineView{ID: ln.ID, Name: ln.Name, Brain: ln.Brain, Repos: repos})
	}
	return views
}

func gaps(cards []RepoCard) []string {
	var out []string
	for _, c := range cards {
		if c.LineID == "" || c.Brain == "à classer" || c.Brain == "" {
			out = append(out, c.FullName)
		}
	}
	sort.Strings(out)
	return out
}

func readGraph(path string) (GraphView, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return GraphView{}, err
	}
	var g struct {
		Nodes []struct {
			ID string `json:"id"`
		} `json:"nodes"`
		Links []struct {
			Source string `json:"source"`
			Target string `json:"target"`
		} `json:"links"`
	}
	if err := json.Unmarshal(raw, &g); err != nil {
		return GraphView{}, err
	}
	repos := map[string]bool{}
	for _, n := range g.Nodes {
		if strings.HasPrefix(n.ID, "ej:") || strings.HasPrefix(n.ID, "ct:") {
			repos[n.ID] = true
		}
	}
	hubs := map[string]bool{"instinct": true, "evolu": true, "cashtro": true}
	bridges := 0
	for _, e := range g.Links {
		if hubs[e.Source] || hubs[e.Target] {
			continue
		}
		if repos[e.Source] && repos[e.Target] {
			bridges++
		}
	}
	return GraphView{RepoNodes: len(repos), Bridges: bridges}, nil
}

func (b Board) counts() (evolu, cashtro int) {
	for _, r := range b.Repos {
		switch r.Org {
		case "Evolu-Jeunes":
			evolu++
		case "cashtro":
			cashtro++
		}
	}
	return evolu, cashtro
}

func findFile(elem ...string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		cand := filepath.Join(append([]string{dir}, elem...)...)
		if _, err := os.Stat(cand); err == nil {
			return cand, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}
