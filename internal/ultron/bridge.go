package ultron

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// BridgeClient talks to the subordinate Cashtro OS.
type BridgeClient struct {
	Base   string
	Client *http.Client
}

func (p *Plane) bridge() *BridgeClient {
	return &BridgeClient{
		Base: strings.TrimRight(p.CashtroURL(), "/"),
		Client: &http.Client{
			Timeout: 8 * time.Second,
		},
	}
}

// PingCashtro hits /health on the subordinate OS.
func (p *Plane) PingCashtro() (map[string]any, bool) {
	b := p.bridge()
	req, err := http.NewRequest(http.MethodGet, b.Base+"/health", nil)
	if err != nil {
		return map[string]any{"error": err.Error()}, false
	}
	res, err := b.Client.Do(req)
	if err != nil {
		return map[string]any{"error": err.Error(), "url": b.Base}, false
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return map[string]any{"status": res.StatusCode, "raw": string(raw)}, res.StatusCode == 200
	}
	data["httpStatus"] = res.StatusCode
	return data, res.StatusCode == 200
}

// BridgeGET proxies a GET to Cashtro.
func (p *Plane) BridgeGET(path string) (int, json.RawMessage, error) {
	b := p.bridge()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	res, err := b.Client.Get(b.Base + path)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	return res.StatusCode, json.RawMessage(raw), err
}

// BridgePOST proxies a POST to Cashtro.
func (p *Plane) BridgePOST(path string, body any) (int, json.RawMessage, error) {
	b := p.bridge()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	var rdr io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		rdr = bytes.NewReader(raw)
	}
	res, err := b.Client.Post(b.Base+path, "application/json", rdr)
	if err != nil {
		return 0, nil, err
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	return res.StatusCode, json.RawMessage(raw), err
}

// InvokeAgent runs a fleet agent — bridges to Cashtro when CashtroID is set.
func (p *Plane) InvokeAgent(actor User, agentID string, payload map[string]any) (map[string]any, error) {
	if !Can(actor.Role, PermFleetInvoke) {
		return nil, ErrForbidden
	}
	p.mu.RLock()
	a, ok := p.agents[agentID]
	p.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}
	if !AllowsCompany(actor, a.CompanyID) {
		return nil, ErrForbidden
	}

	out := map[string]any{
		"agent":     a.ID,
		"name":      a.Name,
		"companyId": a.CompanyID,
		"mode":      a.Mode,
	}

	if a.CashtroID != "" {
		cap := ""
		if len(a.Capabilities) > 0 {
			cap = a.Capabilities[0]
		}
		if v, ok := payload["capability"].(string); ok && v != "" {
			cap = v
		}
		body := map[string]any{"capability": cap, "payload": payload}
		status, raw, err := p.BridgePOST("/api/agents/"+a.CashtroID+"/invoke", body)
		if err != nil {
			p.Audit(actor.Email, "fleet.invoke", agentID, false, err.Error())
			return nil, err
		}
		var data any
		_ = json.Unmarshal(raw, &data)
		out["bridged"] = true
		out["cashtroId"] = a.CashtroID
		out["httpStatus"] = status
		out["result"] = data
		p.Audit(actor.Email, "fleet.invoke", agentID, status >= 200 && status < 300, a.CashtroID)
		p.persist()
		return out, nil
	}

	// Local dry-run when no Cashtro bind.
	out["bridged"] = false
	out["result"] = map[string]any{
		"ok":      true,
		"message": fmt.Sprintf("%s dry-run on Ultron · bind a Cashtro id or worker to execute", a.Name),
		"payload": payload,
	}
	p.Audit(actor.Email, "fleet.invoke", agentID, true, "local-dry-run")
	p.persist()
	return out, nil
}
