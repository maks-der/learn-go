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
	Slug     string
	Order    int
	Title    string
	Excerpt  string
	HTML     template.HTML
	Headings []Heading
	ModTime  time.Time
}

const contentDir = "common"

// Store holds handbook pages from common/en and common/ru.
type Store struct {
	mu     sync.RWMutex
	root   string
	md     goldmark.Markdown
	topics map[string]map[string]Topic // lang -> slug -> topic
	mtime  map[string]time.Time
}

// FindRoot walks from the working directory and the executable
// until it finds common/en and common/ru.
func FindRoot() (string, error) {
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

// NewStore loads all Markdown files from common/en and common/ru.
func NewStore(root string) (*Store, error) {
	s := &Store{
		root: root,
		md: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			goldmark.WithRendererOptions(html.WithHardWraps()),
		),
		topics: map[string]map[string]Topic{
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
	next := map[string]map[string]Topic{
		"en": {},
		"ru": {},
	}
	nextM := map[string]time.Time{}
	for _, lang := range []string{"en", "ru"} {
		dir := filepath.Join(s.root, langDir(lang))
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
				continue
			}
			slug := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
			if !ValidSlug(slug) {
				continue
			}
			path := filepath.Join(dir, e.Name())
			topic, err := s.parseFile(slug, path)
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			next[lang][slug] = topic
			nextM[path] = topic.ModTime
		}
	}
	s.mu.Lock()
	s.topics = next
	s.mtime = nextM
	s.mu.Unlock()
	return nil
}

func (s *Store) parseFile(slug, path string) (Topic, error) {
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
		dir := filepath.Join(s.root, langDir(lang))
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		s.mu.RLock()
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".md") {
				continue
			}
			path := filepath.Join(dir, e.Name())
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
		s.mu.RUnlock()
		if changed {
			break
		}
	}
	if changed {
		_ = s.reload()
	}
}

// List returns topics for a language, sorted by order.
func (s *Store) List(lang string) []Topic {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	src := s.topics[lang]
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

// Get returns one topic by language and slug.
func (s *Store) Get(lang, slug string) (Topic, bool) {
	s.refreshIfChanged()
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.topics[lang][slug]
	return t, ok
}

// Neighbors returns the previous and next topics in the same language.
func (s *Store) Neighbors(lang, slug string) (prev, next *Topic) {
	list := s.List(lang)
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
