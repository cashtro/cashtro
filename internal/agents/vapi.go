package agents

import (
	"fmt"
	"strings"

	"github.com/cashtro/cashtro/internal/kernel"
)

// VapiDesk is Panda's customer service. It files a draft on db:panda.
// A call is not placed. No other line has this desk.
type VapiDesk struct {
	Project  string `json:"project"`
	Database string `json:"database"`
	Role     string `json:"role"`
	Voice    string `json:"voice"`
	Prompt   string `json:"prompt"`
}

// VapiNote is one customer-service draft on Panda's database.
type VapiNote struct {
	Project  string `json:"project"`
	Database string `json:"database"`
	Kind     string `json:"kind"`
	Title    string `json:"title"`
	Body     string `json:"body"`
	Status   string `json:"status"`
	Dialed   bool   `json:"dialed"`
}

// Vapi is the customer service of Panda.
func Vapi() VapiDesk {
	return VapiDesk{
		Project:  "panda",
		Database: "db:panda",
		Role:     "service client",
		Voice:    "vapi",
		Prompt:   "Tu es le service client de Panda. Ta seule base est db:panda. Tu ne lis pas une autre base et tu n'y écris pas. Aucun appel ne part.",
	}
}

// VapiFile stores a customer-service draft for Panda and nowhere else.
func VapiFile(k *kernel.Kernel, project, kind, title, body string) (VapiNote, error) {
	project = strings.TrimSpace(project)
	kind = strings.TrimSpace(kind)
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	desk := Vapi()
	if project != desk.Project {
		return VapiNote{}, fmt.Errorf("Vapi est le service client de Panda")
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
	if strings.Contains(text, "db:") && !strings.Contains(text, desk.Database) {
		return VapiNote{}, fmt.Errorf("autre base refusée")
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
