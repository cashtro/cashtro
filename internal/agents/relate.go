package agents

import (
	"fmt"

	"github.com/cashtro/cashtro/internal/kernel"
)

// ProjectRelation is one proven tie between two repositories.
// Voltron keeps it. It is not a website and it does not deploy.
type ProjectRelation struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Kind   string   `json:"kind"`
	Via    string   `json:"via"`
	Agents []string `json:"agents"`
}

// ProjectRelations is the chain of repositories that share a tree, a clone, or a fiche.
// Each one was read. A name that only sounds alike is not in this list.
func ProjectRelations() []ProjectRelation {
	return []ProjectRelation{
		{From: "Evolu-Jeunes/Pandora", To: "cashtro/PBTM", Kind: "same-tree", Via: "Le même arbre. Le package.json de PBTM pointe Evolu-Jeunes/Pandora.git.", Agents: []string{"binder", "memory"}},
		{From: "Evolu-Jeunes/Panda", To: "Evolu-Jeunes/Pandora", Kind: "same-tree", Via: "Panda et Pandora partagent l'arbre. Panda reste le white-glove.", Agents: []string{"binder", "architect"}},
		{From: "Evolu-Jeunes/CRM", To: "Evolu-Jeunes/CRM-Agents", Kind: "clone", Via: "Même application. Le carnet agents est une autre base.", Agents: []string{"binder", "comms"}},
		{From: "Evolu-Jeunes/CRM", To: "Evolu-Jeunes/Fix2", Kind: "database", Via: "Les clients Fix Tout vivent dans Fix2. Le CRM Lovable ne les reçoit pas.", Agents: []string{"binder", "operator"}},
		{From: "Evolu-Jeunes/H2oH2o", To: "cashtro/H2OriginTest", Kind: "same-tree", Via: "Même site, deux dépôts.", Agents: []string{"binder", "memory"}},
		{From: "Evolu-Jeunes/AI-BOT", To: "cashtro/trading_bot-main", Kind: "same-tree", Via: "Le bot et sa copie cashtro. Aucun ordre live.", Agents: []string{"binder", "security"}},
		{From: "Evolu-Jeunes/sigma", To: "Evolu-Jeunes/sigmaNew", Kind: "same-tree", Via: "sigmaNew reprend sigma.", Agents: []string{"binder", "explorer"}},
		{From: "Evolu-Jeunes/Forge", To: "cashtro/epicenter", Kind: "reports", Via: "Forge rapporte. Epicenter Einstein décide.", Agents: []string{"binder", "manager"}},
	}
}

func binderInvoke(k *kernel.Kernel, call kernel.Call) (kernel.Result, error) {
	if call.Capability != "binder.relate" {
		return kernel.Result{}, fmt.Errorf("%w: %s", kernel.ErrUnknownCapability, call.Capability)
	}
	rels := ProjectRelations()
	k.Publish("binder", "voltron", "relations entre dépôts", map[string]any{"count": len(rels)})
	return kernel.Result{OK: true, Message: "chaîne de relations", Data: rels}, nil
}
