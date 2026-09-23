package agents

import (
	"fmt"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// VapiDesk is one voice agent for one project. Each project has its own database.
// Vapi does not read or write a second database on the same call.
type VapiDesk struct {
	Project  string `json:"project"`
	Database string `json:"database"`
	Voice    string `json:"voice"`
	Shared   bool   `json:"shared"`
	Prompt   string `json:"prompt"`
}

// VapiNote is one draft kept in that project's database only.
type VapiNote struct {
	Project  string `json:"project"`
	Database string `json:"database"`
	Kind     string `json:"kind"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Status   string `json:"status"`
	Dialed   bool   `json:"dialed"`
}

// VapiDesks gives every line its own database. None of them are shared.
func VapiDesks() []VapiDesk {
	lines := Lines()
	out := make([]VapiDesk, 0, len(lines))
	for _, ln := range lines {
		db := vapiDatabase(ln.ID)
		out = append(out, VapiDesk{
			Project:  ln.ID,
			Database: db,
			Voice:    "vapi",
			Shared:   false,
			Prompt:   "Tu es Vapi pour le projet " + ln.ID + " seulement. Ta seule base est " + db + ". Chaque projet a sa propre base. Tu ne lis pas une autre base et tu n'y écris pas.",
		})
	}
	return out
}

// VapiDeskBy returns the voice desk for one project.
func VapiDeskBy(project string) (VapiDesk, bool) {
	project = strings.TrimSpace(project)
	for _, desk := range VapiDesks() {
		if desk.Project == project {
			return desk, true
		}
	}
	return VapiDesk{}, false
}

func vapiDatabase(project string) string {
	if project == "fix2" {
		return "Evolu-Jeunes/Fix2"
	}
	return "db:" + project
}

// VapiFile stores a voice draft in that project's database and nowhere else.
// A call is not placed.
func VapiFile(k *kernel.Kernel, project, kind, title, body string) (VapiNote, error) {
	project = strings.TrimSpace(project)
	kind = strings.TrimSpace(kind)
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	desk, ok := VapiDeskBy(project)
	if !ok {
		return VapiNote{}, fmt.Errorf("projet inconnu")
	}
	if !crmKinds[kind] {
		return VapiNote{}, fmt.Errorf("kind requis: contact, lead, note")
	}
	if title == "" {
		return VapiNote{}, fmt.Errorf("titre requis")
	}
	if names := secretNames(title + "\n" + body); len(names) > 0 {
		return VapiNote{}, fmt.Errorf("secret filtré: %s", strings.Join(names, ", "))
	}
	text := strings.ToLower(title + "\n" + body)
	for _, other := range VapiDesks() {
		if other.Project == project {
			continue
		}
		if strings.Contains(text, strings.ToLower(other.Database)) || strings.Contains(text, "projet "+other.Project) {
			return VapiNote{}, fmt.Errorf("autre base refusée")
		}
	}
	note := VapiNote{
		Project:  desk.Project,
		Database: desk.Database,
		Kind:     kind,
		Title:    title,
		Body:     body,
		Status:   "brouillon",
		Dialed:   false,
	}
	if k != nil {
		k.Remember("vapi:"+desk.Database, kind+" · "+title)
	}
	return note, nil
}
