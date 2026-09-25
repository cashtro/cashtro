package agents

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

// VaultEntry is one record in The Vault. The hash chains it to the record before it.
// A key entry stores the name of the slot. It does not store the secret.
type VaultEntry struct {
	Index   int    `json:"index"`
	Section string `json:"section"`
	Kind    string `json:"kind"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Prev    string `json:"prev"`
	Hash    string `json:"hash"`
}

func (e VaultEntry) seal() string {
	sum := sha256.Sum256([]byte(strconv.Itoa(e.Index) + "|" + e.Prev + "|" + e.Section + "|" + e.Kind + "|" + e.Title + "|" + e.Body))
	return hex.EncodeToString(sum[:])
}

// Vault is the heart. Core memory, money, and keys. More sections can be added later.
type Vault struct {
	Name     string       `json:"name"`
	Heart    bool         `json:"heart"`
	Sections []string     `json:"sections"`
	Entries  []VaultEntry `json:"entries"`
	Brains   []string     `json:"brains"`
	Cells    []string     `json:"cells"`
	Intact   bool         `json:"intact"`
}

// OpenVault builds the library. Books are named. Their pages are not stored.
func OpenVault() Vault {
	v := Vault{
		Name:     "The Vault",
		Heart:    true,
		Sections: []string{"core", "money", "keys"},
	}
	for _, vein := range HustlerVeins() {
		v = v.append("core", "lesson", vein.Source, vein.Lesson)
	}
	for _, book := range bookCatalog() {
		v = v.append("core", "catalog", book[0], book[1])
	}
	v = v.append("money", "ledger", "fonds", "Grand livre interne. Pas une banque. Rien ne sort sans allow.")
	for _, slot := range []string{"OPENROUTER_API_KEY", "DATABASE_URL"} {
		v = v.append("keys", "slot", slot, "")
	}
	for _, b := range Brains() {
		v.Brains = append(v.Brains, b.ID)
	}
	for _, organ := range FusionBody().Organs {
		v.Cells = append(v.Cells, organ.ID)
	}
	v.Intact = vaultIntact(v.Entries)
	return v
}

func bookCatalog() [][2]string {
	return [][2]string{
		{"The 48 Laws of Power", "Robert Greene. Title only."},
		{"The Art of Seduction", "Robert Greene. Title only."},
		{"The 33 Strategies of War", "Robert Greene. Title only."},
		{"Mastery", "Robert Greene. Title only."},
		{"The Laws of Human Nature", "Robert Greene. Title only."},
		{"The 50th Law", "Robert Greene and 50 Cent. Title only."},
		{"The Lean Startup", "Eric Ries. Title only."},
		{"Good to Great", "Jim Collins. Title only."},
		{"Zero to One", "Peter Thiel. Title only."},
		{"Competitive Strategy", "Michael Porter. Title only."},
		{"Profit First", "Mike Michalowicz. Title only."},
		{"The E-Myth Revisited", "Michael Gerber. Title only."},
		{"The Hard Thing About Hard Things", "Ben Horowitz. Title only."},
		{"How to Solve It", "George Polya. Title only."},
		{"The Psychology of Money", "Morgan Housel. Title only."},
		{"Market Wizards", "Jack Schwager. Title only."},
	}
}

func (v Vault) append(section, kind, title, body string) Vault {
	prev := ""
	idx := 0
	if n := len(v.Entries); n > 0 {
		prev = v.Entries[n-1].Hash
		idx = v.Entries[n-1].Index + 1
	}
	entry := VaultEntry{Index: idx, Section: section, Kind: kind, Title: title, Body: body, Prev: prev}
	entry.Hash = entry.seal()
	v.Entries = append(v.Entries, entry)
	v.Intact = vaultIntact(v.Entries)
	return v
}

// Add files one more record. A secret value is refused. A new section is allowed.
func (v Vault) Add(section, kind, title, body string) (Vault, VaultEntry, bool) {
	section = strings.TrimSpace(section)
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if section == "" || title == "" || refusedRule(title+" "+body) {
		return v, VaultEntry{}, false
	}
	if section == "keys" && (body != "" || len(secretNames(title+" "+body)) > 0) {
		return v, VaultEntry{}, false
	}
	if len(secretNames(body)) > 0 {
		return v, VaultEntry{}, false
	}
	known := false
	for _, s := range v.Sections {
		if s == section {
			known = true
		}
	}
	if !known {
		v.Sections = append(v.Sections, section)
	}
	v = v.append(section, kind, title, body)
	return v, v.Entries[len(v.Entries)-1], true
}

func vaultIntact(entries []VaultEntry) bool {
	if len(entries) == 0 {
		return false
	}
	prev := ""
	for i, e := range entries {
		if e.Index != i || e.Prev != prev || e.Hash != e.seal() {
			return false
		}
		prev = e.Hash
	}
	return true
}

// Section returns the records in one room of the vault.
func (v Vault) Section(name string) []VaultEntry {
	var out []VaultEntry
	for _, e := range v.Entries {
		if e.Section == name {
			out = append(out, e)
		}
	}
	return out
}
