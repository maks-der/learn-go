package content

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
)

var (
	titleLine = regexp.MustCompile(`(?m)^#\s+(.+)$`)
	slugSafe  = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	descHeads = []string{"Description", "Описание"}
)

// Heading is a level-2 section used for in-page contents.
type Heading struct {
	ID   string
	Text string
}

// Topic is one handbook page loaded from a Markdown file.
type Topic struct {
	Section  string
	Slug     string
	Order    int
	Title    string
	Excerpt  string
	HTML     template.HTML
	Headings []Heading
	ModTime  time.Time
}

const contentDir = "common"

// Store holds handbook pages from common/{lang}/{section}.
type Store struct {
	mu     sync.RWMutex
	root   string
	md     goldmark.Markdown
	topics map[string]map[string]map[string]Topic // lang -> section -> slug
	mtime  map[string]time.Time
}

// FindRoot walks from the working directory and the executable
// until it finds common/en and common/ru.
// Set LEARNGO_ROOT to skip the search (used on hosts such as Render).
func FindRoot() (string, error) {
	if root := os.Getenv("LEARNGO_ROOT"); root != "" {
		if isContentRoot(root) {
			return root, nil
		}
		return "", fmt.Errorf("LEARNGO_ROOT is not a content root: %s", root)
	}
	seen := map[string]bool{}
	var starts []string
	if cwd, err := os.Getwd(); err == nil {
		starts = append(starts, cwd)
	}
	if exe, err := os.Executable(); err == nil {
		starts = append(starts, filepath.Dir(exe))
	}
	for _, start := range starts {
		for dir := start; ; dir = filepath.Dir(dir) {
			if seen[dir] {
				break
			}
			seen[dir] = true
			if isContentRoot(dir) {
				return dir, nil
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	return "", fmt.Errorf("common/en and common/ru folders not found")
}

func isContentRoot(dir string) bool {
	en, err1 := os.Stat(filepath.Join(dir, contentDir, "en"))
	ru, err2 := os.Stat(filepath.Join(dir, contentDir, "ru"))
	return err1 == nil && err2 == nil && en.IsDir() && ru.IsDir()
}

// NewStore loads Markdown files from common/{lang}/{section}.
func NewStore(root string) (*Store, error) {
	s := &Store{
		root: root,
		md: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(html.WithHardWraps()),
		),
		topics: map[string]map[string]map[string]Topic{
			"en": {},
			"ru": {},
		},
		mtime: map[string]time.Time{},
	}
	if err := s.reload(); err != nil {
		return nil, err
	}
	return s, nil
}

func langDir(lang string) string {
	switch lang {
	case "ru":
		return filepath.Join(contentDir, "ru")
	default:
		return filepath.Join(contentDir, "en")
	}
}

// ValidLang reports whether lang is en or ru.
func ValidLang(lang string) bool {
	return lang == "en" || lang == "ru"
}

// ValidSlug reports whether slug is a safe file stem.
func ValidSlug(slug string) bool {
	return slugSafe.MatchString(slug)
}

func (s *Store) reload() error {
	next := map[string]map[string]map[string]Topic{
		"en": {},
		"ru": {},
	}
	nextM := map[string]time.Time{}
	for _, lang := range []string{"en", "ru"} {
		langPath := filepath.Join(s.root, langDir(lang))
		entries, err := os.ReadDir(langPath)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if !e.IsDir() || !ValidSection(e.Name()) {
				continue
			}
			section := e.Name()
			secPath := filepath.Join(langPath, section)
			files, err := os.ReadDir(secPath)
			if err != nil {
				return err
			}
			if next[lang][section] == nil {
				next[lang][section] = map[string]Topic{}
			}
			for _, f := range files {
				if f.IsDir() || !strings.HasSuffix(strings.ToLower(f.Name()), ".md") {
					continue
				}
				slug := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
				if !ValidSlug(slug) {
					continue
				}
				path := filepath.Join(secPath, f.Name())
				topic, err := s.parseFile(section, slug, path)
				if err != nil {
					return fmt.Errorf("%s: %w", path, err)
				}
				next[lang][section][slug] = topic
				nextM[path] = topic.ModTime
			}
		}
	}
	s.mu.Lock()
	s.topics = next
	s.mtime = nextM
	s.mu.Unlock()
	return nil
}

func (s *Store) parseFile(section, slug, path string) (Topic, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Topic{}, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return Topic{}, err
	}
	src := text.NewReader(raw)
	doc := s.md.Parser().Parse(src)

	var buf bytes.Buffer
	if err := s.md.Renderer().Render(&buf, raw, doc); err != nil {
		return Topic{}, err
	}

	title := slug
	if m := titleLine.FindSubmatch(raw); len(m) == 2 {
		title = strings.TrimSpace(string(m[1]))
	}

	return Topic{
		Section:  section,
		Slug:     slug,
		Order:    orderFromSlug(slug),
		Title:    title,
		Excerpt:  excerptFrom(raw),
		HTML:     template.HTML(buf.String()),
		Headings: collectH2(doc, raw),
		ModTime:  info.ModTime(),
	}, nil
}

func orderFromSlug(slug string) int {
	i := strings.IndexByte(slug, '-')
	if i <= 0 {
		return 0
	}
	n, err := strconv.Atoi(slug[:i])
	if err != nil {
		return 0
	}
	return n
}

func excerptFrom(raw []byte) string {
	text := string(raw)
	lower := strings.ToLower(text)
	start := -1
	for _, h := range descHeads {
		key := "## " + strings.ToLower(h)
		idx := strings.Index(lower, key)
		if idx >= 0 {
			start = idx + len(key)
			break
		}
	}
	if start < 0 {
		return ""
	}
	rest := text[start:]
	rest = strings.TrimLeft(rest, "\r\n \t")
	var b strings.Builder
	for _, line := range strings.Split(rest, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "---" {
			break
		}
		if line == "" {
			if b.Len() > 0 {
				break
			}
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(line)
	}
	out := strings.TrimSpace(b.String())
	if len([]rune(out)) > 220 {
		out = string([]rune(out)[:220]) + "…"
	}
	return out
}

func collectH2(doc ast.Node, raw []byte) []Heading {
	var out []Heading
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := n.(*ast.Heading)
		if !ok || h.Level != 2 {
			return ast.WalkContinue, nil
		}
		text := strings.TrimSpace(plainText(h, raw))
		if text == "" {
			return ast.WalkContinue, nil
		}
		id := ""
		if v, ok := h.AttributeString("id"); ok {
			if s, ok := v.([]byte); ok {
				id = string(s)
			} else if s, ok := v.(string); ok {
				id = s
			}
		}
		out = append(out, Heading{ID: id, Text: text})
		return ast.WalkContinue, nil
	})
	return out
}

func plainText(n ast.Node, raw []byte) string {
	var b strings.Builder
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		if t, ok := c.(*ast.Text); ok {
			b.Write(t.Segment.Value(raw))
		}
		return ast.WalkContinue, nil
	})
	return b.String()
}

func (s *Store) refreshIfChanged() {
	changed := false
	for _, lang := range []string{"en", "ru"} {
		langPath := filepath.Join(s.root, langDir(lang))
		entries, err := os.ReadDir(langPath)
		if err != nil {
			return
		}
		s.mu.RLock()
		for _, e := range entries {
			if !e.IsDir() || !ValidSection(e.Name()) {
				continue
			}
			secPath := filepath.Join(langPath, e.Name())
			files, err := os.ReadDir(secPath)
			if err != nil {
				continue
			}
			for _, f := range files {
				if f.IsDir() || !strings.HasSuffix(strings.ToLower(f.Name()), ".md") {
					continue
				}
				path := filepath.Join(secPath, f.Name())
				info, err := os.Stat(path)
				if err != nil {
					continue
				}
				prev, ok := s.mtime[path]
				if !ok || !info.ModTime().Equal(prev) {
					changed = true
					break
				}
			}
			if changed {
				break
			}
		}
		s.mu.RUnlock()
		if changed {
			break
		}
	}
	if changed {
		_ = s.reload()
	}
}

// HasSection reports whether a section exists for the language.
func (s *Store) HasSection(lang, section string) bool {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.topics[lang][section]
	return ok
}

// ListSections returns subject folders for a language.
func (s *Store) ListSections(lang string) []Section {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.topics[lang]
	out := make([]Section, 0, len(src))
	for slug, topics := range src {
		out = append(out, infoForSection(lang, slug, len(topics)))
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Order == out[j].Order {
			return out[i].Slug < out[j].Slug
		}
		return out[i].Order < out[j].Order
	})
	return out
}

// Section returns metadata for one section.
func (s *Store) Section(lang, section string) (Section, bool) {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	topics, ok := s.topics[lang][section]
	if !ok {
		return Section{}, false
	}
	return infoForSection(lang, section, len(topics)), true
}

// List returns topics for a language and section, sorted by order.
func (s *Store) List(lang, section string) []Topic {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.topics[lang][section]
	out := make([]Topic, 0, len(src))
	for _, t := range src {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Order == out[j].Order {
			return out[i].Slug < out[j].Slug
		}
		return out[i].Order < out[j].Order
	})
	return out
}

// Get returns one topic by language, section, and slug.
func (s *Store) Get(lang, section, slug string) (Topic, bool) {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.topics[lang][section][slug]
	return t, ok
}

// FindSlug looks up a topic slug in any section for redirects.
func (s *Store) FindSlug(lang, slug string) (Topic, bool) {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, topics := range s.topics[lang] {
		if t, ok := topics[slug]; ok {
			return t, true
		}
	}
	return Topic{}, false
}

// Neighbors returns the previous and next topics in the same section.
func (s *Store) Neighbors(lang, section, slug string) (prev, next *Topic) {
	list := s.List(lang, section)
	for i, t := range list {
		if t.Slug != slug {
			continue
		}
		if i > 0 {
			p := list[i-1]
			prev = &p
		}
		if i+1 < len(list) {
			n := list[i+1]
			next = &n
		}
		break
	}
	return prev, next
}
