package team

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func sampleConfig() Config {
	return Config{
		Project:     "MyProject",
		TmuxSession: "myproject",
		Board:       "board/BOARD.md",
		Roles:       []Role{{Name: "lead", Folder: "lead"}, {Name: "eng", Folder: "eng"}},
	}
}

func TestConfigWriteAndDiscover(t *testing.T) {
	root := t.TempDir()
	if err := sampleConfig().WriteTo(root); err != nil {
		t.Fatalf("WriteTo: %v", err)
	}
	// Discover from a nested dir walks up to the project root.
	start := filepath.Join(root, "eng", "src")
	if err := os.MkdirAll(start, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg, gotRoot, ok := Discover(start)
	if !ok {
		t.Fatal("Discover: expected to find .team/config.json")
	}
	if gotRoot != root {
		t.Errorf("root = %q, want %q", gotRoot, root)
	}
	if cfg.TmuxSession != "myproject" || len(cfg.Roles) != 2 {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestDiscoverNoConfigIsNoOp(t *testing.T) {
	if _, _, ok := Discover(t.TempDir()); ok {
		t.Error("Discover should return ok=false when no .team/config.json exists")
	}
}

func TestRoleForDir(t *testing.T) {
	root := t.TempDir()
	cfg := sampleConfig()
	got, ok := RoleForDir(cfg, root, filepath.Join(root, "eng", "deep", "path"))
	if !ok || got.Name != "eng" {
		t.Errorf("RoleForDir(eng subdir) = %+v, %v; want eng", got, ok)
	}
	if _, ok := RoleForDir(cfg, root, filepath.Join(root, "unknown")); ok {
		t.Error("RoleForDir(unknown) should be false")
	}
}

func TestNotifyTargetsOnlyActiveRequests(t *testing.T) {
	cfg := sampleConfig()
	board := strings.Join([]string{
		"## Cross-Role Requests",
		"- @eng: build the login API. DoD: test + curl.",
		"- @lead: [DONE 2026-08-11] set the roadmap.", // marked done → no poke
	}, "\n")
	got := NotifyTargets(cfg, board)
	if len(got) != 1 || got[0] != "myproject:eng" {
		t.Errorf("targets = %v; want [myproject:eng]", got)
	}
}

func TestFreshBoardDoesNotSelfTrigger(t *testing.T) {
	// A freshly rendered board (with its example line) must produce no
	// notify targets and no pending requests — the example must not use a
	// real role name.
	roles := []RoleSpec{{Name: "lead", Folder: "lead"}, {Name: "eng", Folder: "eng"}}
	board := RenderBoard("Proj", "ship it", roles)
	if got := NotifyTargets(sampleConfig(), board); len(got) != 0 {
		t.Errorf("fresh board self-triggers notify: %v", got)
	}
	if got := PendingRequests(board, "eng"); len(got) != 0 {
		t.Errorf("fresh board has phantom pending eng requests: %v", got)
	}
}

func TestPendingRequestsSkipsDone(t *testing.T) {
	board := strings.Join([]string{
		"- @eng: task one",
		"- @eng: [DONE] task two",
		"- @marketing: not mine",
	}, "\n")
	got := PendingRequests(board, "eng")
	if len(got) != 1 || got[0] != "- @eng: task one" {
		t.Errorf("pending = %v; want one active eng line", got)
	}
}

func TestSlug(t *testing.T) {
	cases := map[string]string{
		"Eng":              "eng",
		"Growth Marketing": "growth-marketing",
		"  R&D  ":          "r-d",
		"data/ops":         "data-ops",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q; want %q", in, got, want)
		}
	}
}

func TestExtractJSONObject(t *testing.T) {
	wrapped := "Here you go:\n```json\n{\"board_focus\":\"x\",\"roles\":[]}\n```\nthanks"
	got := extractJSONObject(wrapped)
	if got != `{"board_focus":"x","roles":[]}` {
		t.Errorf("extractJSONObject = %q", got)
	}
	// Braces inside strings must not end the object early.
	tricky := `{"a":"}","b":1}`
	if got := extractJSONObject(tricky); got != tricky {
		t.Errorf("extractJSONObject(tricky) = %q; want %q", got, tricky)
	}
}

func TestEnsureHooksIdempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	// Pre-existing unrelated content must be preserved.
	if err := os.WriteFile(path, []byte(`{"model":"opus","hooks":{"SessionStart":[]}}`), 0o644); err != nil {
		t.Fatal(err)
	}
	changed, err := EnsureHooks(path)
	if err != nil || !changed {
		t.Fatalf("first EnsureHooks: changed=%v err=%v", changed, err)
	}
	// Second run must be a no-op (no duplicate hooks).
	changed2, err := EnsureHooks(path)
	if err != nil || changed2 {
		t.Fatalf("second EnsureHooks should be no-op: changed=%v err=%v", changed2, err)
	}
	data, _ := os.ReadFile(path)
	s := string(data)
	if !strings.Contains(s, `"model": "opus"`) {
		t.Error("EnsureHooks clobbered unrelated keys")
	}
	if n := strings.Count(s, "ck team hook notify"); n != 1 {
		t.Errorf("notify hook appears %d times; want 1", n)
	}
	if !strings.Contains(s, "ck team hook rescan") {
		t.Error("rescan hook missing")
	}
}
