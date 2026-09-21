package ultron

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	sessionTTL   = 12 * time.Hour
	hashIters    = 120000
	saltBytes    = 16
	tokenBytes   = 32
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadLogin     = errors.New("invalid credentials")
	ErrExists       = errors.New("already exists")
	ErrNotFound     = errors.New("not found")
	ErrInvalid      = errors.New("invalid")
)

func randomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

func hashPassword(password, saltB64 string) (string, error) {
	salt, err := base64.RawStdEncoding.DecodeString(saltB64)
	if err != nil {
		return "", err
	}
	sum := []byte(password)
	for i := 0; i < hashIters; i++ {
		h := sha256.New()
		h.Write(salt)
		h.Write(sum)
		sum = h.Sum(nil)
	}
	return base64.RawStdEncoding.EncodeToString(sum), nil
}

func newSalt() (string, error) {
	b, err := randomBytes(saltBytes)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}

func newToken() (string, error) {
	b, err := randomBytes(tokenBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// SetPassword hashes and stores a password on the user.
func (p *Plane) SetPassword(u *User, password string) error {
	if len(password) < 8 {
		return fmt.Errorf("%w: password too short", ErrInvalid)
	}
	salt, err := newSalt()
	if err != nil {
		return err
	}
	hash, err := hashPassword(password, salt)
	if err != nil {
		return err
	}
	u.Salt = salt
	u.PasswordHash = hash
	return nil
}

func (p *Plane) checkPassword(u *User, password string) bool {
	got, err := hashPassword(password, u.Salt)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(u.PasswordHash)) == 1
}

// Login authenticates and returns a session token.
func (p *Plane) Login(email, password string) (Session, User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	p.mu.Lock()
	var u *User
	for _, cand := range p.users {
		if strings.EqualFold(cand.Email, email) {
			u = cand
			break
		}
	}
	if u == nil || !u.Active || !p.checkPassword(u, password) {
		p.mu.Unlock()
		p.Audit(email, "login", "", false, "bad credentials")
		return Session{}, User{}, ErrBadLogin
	}
	tok, err := newToken()
	if err != nil {
		p.mu.Unlock()
		return Session{}, User{}, err
	}
	now := p.now()
	s := &Session{Token: tok, UserID: u.ID, CreatedAt: now, ExpiresAt: now.Add(sessionTTL)}
	p.sessions[tok] = s
	outU := u.Public()
	outS := *s
	p.mu.Unlock()
	p.Audit(outU.Email, "login", "", true, "ok")
	p.persist()
	return outS, outU, nil
}

// Logout drops a session.
func (p *Plane) Logout(token string) {
	p.mu.Lock()
	if s, ok := p.sessions[token]; ok {
		delete(p.sessions, token)
		p.mu.Unlock()
		p.Audit(s.UserID, "logout", "", true, "")
		p.persist()
		return
	}
	p.mu.Unlock()
}

// UserForToken resolves a bearer token.
func (p *Plane) UserForToken(token string) (User, bool) {
	token = strings.TrimSpace(token)
	if token == "" {
		return User{}, false
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	s, ok := p.sessions[token]
	if !ok {
		return User{}, false
	}
	if p.now().After(s.ExpiresAt) {
		delete(p.sessions, token)
		return User{}, false
	}
	u, ok := p.users[s.UserID]
	if !ok || !u.Active {
		return User{}, false
	}
	return u.Public(), true
}

// CreateUser adds a principal (owner/admin only at call site).
func (p *Plane) CreateUser(actor User, email, name string, role Role, companyIDs []string, password string) (User, error) {
	if !Can(actor.Role, PermUsersWrite) {
		return User{}, ErrForbidden
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || name == "" {
		return User{}, ErrInvalid
	}
	if role == "" {
		role = RoleViewer
	}
	p.mu.Lock()
	for _, u := range p.users {
		if strings.EqualFold(u.Email, email) {
			p.mu.Unlock()
			return User{}, ErrExists
		}
	}
	id := slugify(strings.Split(email, "@")[0])
	if id == "" {
		id = "user"
	}
	base := id
	for n := 2; ; n++ {
		if _, ok := p.users[id]; !ok {
			break
		}
		id = fmt.Sprintf("%s-%d", base, n)
	}
	u := &User{
		ID:         id,
		Email:      email,
		Name:       strings.TrimSpace(name),
		Role:       role,
		CompanyIDs: append([]string(nil), companyIDs...),
		Active:     true,
		CreatedAt:  p.now(),
	}
	if err := p.SetPassword(u, password); err != nil {
		p.mu.Unlock()
		return User{}, err
	}
	p.users[id] = u
	out := u.Public()
	p.mu.Unlock()
	p.Audit(actor.Email, "user.create", id, true, string(role))
	p.persist()
	return out, nil
}

// ListUsers returns public users.
func (p *Plane) ListUsers() []User {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]User, 0, len(p.users))
	for _, u := range p.users {
		out = append(out, u.Public())
	}
	return out
}
