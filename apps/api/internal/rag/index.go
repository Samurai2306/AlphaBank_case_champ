package rag

import (
	"bufio"
	"bytes"
	"encoding/json"
	"embed"
	"strings"
	"sync"
	"unicode"
)

//go:embed data/*.jsonl
var corpusFS embed.FS

var (
	defaultOnce sync.Once
	defaultIdx  *Index
)

// Index is an in-memory keyword index over the demo KB corpus.
type Index struct {
	chunks []Chunk
}

// Default returns the process-wide corpus index (loaded once).
func Default() *Index {
	defaultOnce.Do(func() {
		idx, err := LoadEmbedded()
		if err != nil {
			defaultIdx = &Index{}
			return
		}
		defaultIdx = idx
	})
	return defaultIdx
}

// LoadEmbedded reads all embedded *.jsonl corpus files.
func LoadEmbedded() (*Index, error) {
	entries, err := corpusFS.ReadDir("data")
	if err != nil {
		return nil, err
	}
	var chunks []Chunk
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		raw, err := corpusFS.ReadFile("data/" + e.Name())
		if err != nil {
			return nil, err
		}
		sc := bufio.NewScanner(bytes.NewReader(raw))
		// Long FAQ lines can exceed default 64K — raise buffer.
		buf := make([]byte, 0, 64*1024)
		sc.Buffer(buf, 1024*1024)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			var c Chunk
			if err := json.Unmarshal([]byte(line), &c); err != nil {
				return nil, err
			}
			if c.ID == "" || c.Text == "" {
				continue
			}
			chunks = append(chunks, c)
		}
		if err := sc.Err(); err != nil {
			return nil, err
		}
	}
	return &Index{chunks: chunks}, nil
}

// Len returns corpus size.
func (idx *Index) Len() int {
	if idx == nil {
		return 0
	}
	return len(idx.chunks)
}

// Retrieve returns topK chunks ranked by lexical overlap with the query.
func (idx *Index) Retrieve(query string, topK int) []Hit {
	if idx == nil || len(idx.chunks) == 0 {
		return nil
	}
	if topK <= 0 {
		topK = 4
	}
	qTokens := tokenize(query)
	if len(qTokens) == 0 {
		return nil
	}
	type scored struct {
		hit Hit
	}
	var ranked []scored
	for _, c := range idx.chunks {
		score := scoreChunk(c, qTokens, query)
		if score < 1.2 {
			continue
		}
		ranked = append(ranked, scored{Hit{
			ID: c.ID, Title: c.Title, Text: c.Text, Source: c.Source, Score: score,
		}})
	}
	// Simple selection sort for tiny N.
	for i := 0; i < len(ranked); i++ {
		best := i
		for j := i + 1; j < len(ranked); j++ {
			if ranked[j].hit.Score > ranked[best].hit.Score {
				best = j
			}
		}
		ranked[i], ranked[best] = ranked[best], ranked[i]
	}
	if len(ranked) > topK {
		ranked = ranked[:topK]
	}
	out := make([]Hit, 0, len(ranked))
	for _, r := range ranked {
		out = append(out, r.hit)
	}
	return out
}

// LowConfidence reports whether hits are too weak for grounded answers.
func LowConfidence(hits []Hit) bool {
	if len(hits) == 0 {
		return true
	}
	return hits[0].Score < 2.5
}

func scoreChunk(c Chunk, qTokens []string, rawQuery string) float64 {
	q := normalize(rawQuery)
	title := normalize(c.Title)
	text := normalize(c.Text)
	var score float64
	for _, tag := range c.Tags {
		nt := normalize(tag)
		if nt == "" {
			continue
		}
		if strings.Contains(q, nt) || tokenIn(qTokens, nt) {
			score += 3.2
		}
	}
	for _, t := range qTokens {
		if len([]rune(t)) < 3 {
			continue
		}
		if strings.Contains(title, t) {
			score += 2.0
		}
		if strings.Contains(text, t) {
			score += 1.0
		}
		if strings.Contains(normalize(c.ID), t) {
			score += 0.5
		}
	}
	// Phrase bonuses for common demo questions.
	phrases := []struct {
		p string
		w float64
	}{
		{"что такое нпд", 4},
		{"что такое усн", 4},
		{"когда платить", 3},
		{"кэшбек", 3},
		{"кешбек", 3},
		{"копил", 2.5},
		{"тариф", 2.5},
		{"рко", 3},
		{"115", 3},
		{"светофор", 3},
	}
	for _, ph := range phrases {
		if strings.Contains(q, ph.p) && (strings.Contains(text, ph.p) || strings.Contains(title, ph.p) || tagsContain(c.Tags, ph.p)) {
			score += ph.w
		}
	}
	return score
}

func tagsContain(tags []string, needle string) bool {
	for _, t := range tags {
		if strings.Contains(normalize(t), needle) {
			return true
		}
	}
	return false
}

func tokenIn(tokens []string, phrase string) bool {
	pt := tokenize(phrase)
	if len(pt) == 1 {
		for _, t := range tokens {
			if t == pt[0] {
				return true
			}
		}
	}
	return false
}

func tokenize(s string) []string {
	s = normalize(s)
	var b strings.Builder
	var out []string
	flush := func() {
		if b.Len() == 0 {
			return
		}
		tok := b.String()
		b.Reset()
		if len([]rune(tok)) >= 2 {
			out = append(out, tok)
		}
	}
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		} else {
			flush()
		}
	}
	flush()
	return out
}

func normalize(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	return strings.ReplaceAll(s, "ё", "е")
}
