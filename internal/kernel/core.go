package kernel

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// KernelLink is one kernel in the chain. It has its own agents and infrastructure.
type KernelLink struct {
	Index          int      `json:"index"`
	ID             string   `json:"id"`
	Parent         string   `json:"parent"`
	Infrastructure string   `json:"infrastructure"`
	Agents         []string `json:"agents"`
	Prev           string   `json:"prev"`
	Hash           string   `json:"hash"`
}

func (l KernelLink) seal() string {
	sum := sha256.Sum256([]byte(strconv.Itoa(l.Index) + "|" + l.Prev + "|" + l.ID + "|" + l.Parent + "|" + l.Infrastructure + "|" + strings.Join(l.Agents, ",")))
	return hex.EncodeToString(sum[:])
}

// Core creates kernels. Each kernel keeps its own agents and its own infrastructure.
type Core struct {
	mu    sync.Mutex
	root  *Kernel
	chain []KernelLink
	kids  map[string]*Kernel
}

// Core returns the chain maker for this kernel.
func (k *Kernel) Core() *Core {
	k.mu.Lock()
	defer k.mu.Unlock()
	if k.core == nil {
		k.core = &Core{root: k, kids: map[string]*Kernel{}}
	}
	return k.core
}

// Chain is the kernel chain, oldest first.
func (c *Core) Chain() []KernelLink {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]KernelLink, len(c.chain))
	copy(out, c.chain)
	return out
}

// Kernel returns one kernel in the chain.
func (c *Core) Kernel(id string) (*Kernel, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	k, ok := c.kids[id]
	return k, ok
}

// SpawnKernel creates the next kernel. Its agents are not the parent's agents.
func (c *Core) SpawnKernel(ctx context.Context, id, infrastructure string, specs []Spec) (KernelLink, error) {
	id = strings.TrimSpace(id)
	infrastructure = strings.TrimSpace(infrastructure)
	if id == "" || infrastructure == "" {
		return KernelLink{}, fmt.Errorf("id and infrastructure required")
	}
	c.mu.Lock()
	if _, ok := c.kids[id]; ok {
		c.mu.Unlock()
		return KernelLink{}, fmt.Errorf("kernel exists: %s", id)
	}
	parent := "root"
	prev := ""
	index := 0
	if n := len(c.chain); n > 0 {
		parent = c.chain[n-1].ID
		prev = c.chain[n-1].Hash
		index = c.chain[n-1].Index + 1
	}
	c.mu.Unlock()

	kid := New()
	if len(specs) == 0 {
		specs = []Spec{{
			ID: id + "-agent", Name: id, Kind: KindUser, Mode: ModeLive,
			Role: "agent", Summary: "Created by the core for " + id,
			Capabilities: []string{id + ".run"}, Autostart: true,
		}}
	}
	names := make([]string, 0, len(specs))
	seen := map[string]bool{}
	for _, spec := range specs {
		spec.ID = strings.TrimSpace(spec.ID)
		spec.Name = strings.TrimSpace(spec.Name)
		if spec.ID == "" || spec.Name == "" || seen[spec.ID] {
			return KernelLink{}, fmt.Errorf("each agent needs its own id and name")
		}
		seen[spec.ID] = true
		if spec.Kind == "" {
			spec.Kind = KindUser
		}
		if spec.Mode == "" {
			spec.Mode = ModeLive
		}
		if len(spec.Capabilities) == 0 {
			spec.Capabilities = []string{spec.ID + ".run"}
		}
		spec.Autostart = true
		if _, err := kid.CreateAgent(ctx, spec); err != nil {
			return KernelLink{}, err
		}
		names = append(names, spec.ID)
	}
	link := KernelLink{Index: index, ID: id, Parent: parent, Infrastructure: infrastructure, Agents: names, Prev: prev}
	link.Hash = link.seal()

	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.kids[id]; ok {
		return KernelLink{}, fmt.Errorf("kernel exists: %s", id)
	}
	c.kids[id] = kid
	c.chain = append(c.chain, link)
	return link, nil
}

// CreateAgent registers a new agent on this kernel and starts it.
func (k *Kernel) CreateAgent(ctx context.Context, spec Spec) (Process, error) {
	spec.ID = strings.TrimSpace(spec.ID)
	spec.Name = strings.TrimSpace(spec.Name)
	if spec.ID == "" || spec.Name == "" {
		return Process{}, fmt.Errorf("agent id and name required")
	}
	k.mu.Lock()
	if _, ok := k.agents[spec.ID]; ok {
		k.mu.Unlock()
		return Process{}, fmt.Errorf("agent exists: %s", spec.ID)
	}
	k.mu.Unlock()
	k.Register(&madeAgent{spec: spec})
	if err := k.Spawn(ctx, spec.ID); err != nil {
		return Process{}, err
	}
	return k.Process(spec.ID)
}

type madeAgent struct {
	spec Spec
}

func (a *madeAgent) Spec() Spec { return a.spec }

func (a *madeAgent) Boot(context.Context, *Kernel) error { return nil }

func (a *madeAgent) Invoke(_ context.Context, call Call) (Result, error) {
	return Result{OK: true, Message: a.spec.Name + " running", Data: map[string]any{
		"id": a.spec.ID, "capability": call.Capability,
	}}, nil
}
