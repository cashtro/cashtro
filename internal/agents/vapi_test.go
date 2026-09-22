package agents

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/cashtro/cashtro/internal/kernel"
	"github.com/cashtro/cashtro/internal/vapi"
)

func TestVapiLaneTalkVoiceCallAndBridges(t *testing.T) {
	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}

	res, err := k.Invoke(context.Background(), "vapi", kernel.Call{Capability: "vapi.status"})
	if err != nil || !res.OK {
		t.Fatalf("status: %+v %v", res, err)
	}
	st := res.Data.(map[string]any)
	if st["product"] != vapi.ProductURL {
		t.Fatalf("status product = %#v", st["product"])
	}
	bridges, _ := st["bridges"].([]string)
	if len(bridges) != 10 {
		t.Fatalf("bridges = %v", bridges)
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{Capability: "vapi.voices"})
	if err != nil || !res.OK {
		t.Fatalf("voices: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.voice",
		Payload:    []byte(`{"id":"denise"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("voice: %+v %v", res, err)
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.bridge",
		Payload:    []byte(`{"line":"proximity"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("bridge: %+v %v", res, err)
	}
	if !strings.Contains(res.Message, "Proximity") {
		t.Fatalf("bridge msg = %q", res.Message)
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.talk",
		Payload:    []byte(`{"text":"where can I change the voice"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("talk: %+v %v", res, err)
	}
	turn := res.Data.(TalkTurn)
	if turn.Live {
		t.Fatal("unbound talk should stay local")
	}
	if !strings.Contains(strings.ToLower(turn.Reply), "voice") {
		t.Fatalf("reply = %q", turn.Reply)
	}
	if turn.Line != "proximity" {
		t.Fatalf("talk line = %q", turn.Line)
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.call",
		Payload:    []byte(`{"to":"+15555550100","prompt":"intern standup","line":"panda"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("call: %+v %v", res, err)
	}
	confirms := k.Confirms()
	if len(confirms) == 0 || confirms[len(confirms)-1].Cap != "vapi.call" {
		t.Fatalf("expected parked vapi.call, got %#v", confirms)
	}
	cid := confirms[len(confirms)-1].ID

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.fire",
		Payload:    []byte(`{"id":` + strconv.Itoa(cid) + `}`),
	})
	if err != nil || res.OK {
		t.Fatalf("fire before allow should fail: %+v %v", res, err)
	}

	if _, err := k.DecideConfirm(cid, true); err != nil {
		t.Fatal(err)
	}
	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.fire",
		Payload:    []byte(`{"confirmId":` + strconv.Itoa(cid) + `}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("fire after allow: %+v %v", res, err)
	}
	if !strings.Contains(res.Message, "queued") && !strings.Contains(res.Message, "bind") {
		t.Fatalf("unbound fire = %q", res.Message)
	}

	res, err = k.Invoke(context.Background(), "vapi", kernel.Call{Capability: "vapi.web"})
	if err != nil || !res.OK {
		t.Fatalf("web: %+v %v", res, err)
	}

	for _, ln := range Links() {
		has := false
		for _, id := range ln.Agents {
			if id == "vapi" {
				has = true
				break
			}
		}
		if !has {
			t.Fatalf("link %s→%s missing vapi", ln.From, ln.To)
		}
	}
}

func TestVapiLiveTalkWithoutAssistantID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret" {
			http.Error(w, "auth", http.StatusUnauthorized)
			return
		}
		switch {
		case r.Method == "GET" && r.URL.Path == "/assistant":
			_, _ = w.Write([]byte(`[]`))
		case r.Method == "POST" && r.URL.Path == "/assistant":
			_, _ = w.Write([]byte(`{"id":"asst_live","name":"Cashtro Teal"}`))
		case r.Method == "POST" && r.URL.Path == "/chat":
			_, _ = w.Write([]byte(`{"id":"chat-9","output":[{"role":"assistant","content":"hey Castro from Vapi"}]}`))
		default:
			http.Error(w, r.URL.Path, http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)
	t.Setenv("VAPI_API_KEY", "secret")
	t.Setenv("VAPI_BASE_URL", srv.URL)
	t.Setenv("VAPI_ASSISTANT_ID", "")

	k, err := Boot()
	if err != nil {
		t.Fatal(err)
	}
	res, err := k.Invoke(context.Background(), "vapi", kernel.Call{
		Capability: "vapi.talk",
		Payload:    []byte(`{"text":"hello Castro"}`),
	})
	if err != nil || !res.OK {
		t.Fatalf("talk: %+v %v", res, err)
	}
	turn := res.Data.(TalkTurn)
	if !turn.Live {
		t.Fatalf("expected live Vapi talk, got %+v", turn)
	}
	if turn.Reply != "hey Castro from Vapi" {
		t.Fatalf("reply = %q", turn.Reply)
	}
}
