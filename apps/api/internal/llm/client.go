package llm

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Browser-like UA: Cloudflare on some free gateways blocks Go's default User-Agent (1010).
const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"

const openCodeZenUA = "opencode/1.15.0 ai-sdk/provider-utils/4.0.23 runtime/bun/1.3.13"

// Endpoint is one OpenAI-compatible chat completions backend.
type Endpoint struct {
	Name     string
	BaseURL  string
	APIKey   string
	Model    string
	Provider string // "", "openai", "opencode-zen"
}

// Client talks to OpenAI-compatible Chat Completions APIs with failover.
type Client struct {
	HTTP       *http.Client
	UserAgent  string
	mu         sync.Mutex
	eps        []Endpoint
	failUntil  map[string]time.Time
	lastOK     string
	lastErr    string
	lastCheck  time.Time
}

func New(baseURL, apiKey, model string) *Client {
	c := &Client{
		HTTP:      &http.Client{Timeout: 22 * time.Second},
		UserAgent: defaultUserAgent,
		failUntil: map[string]time.Time{},
	}
	c.SetEndpoints([]Endpoint{{
		Name: "primary", BaseURL: baseURL, APIKey: apiKey, Model: model,
	}})
	return c
}

// NewWithFailover builds a client that rotates through endpoints on errors.
func NewWithFailover(eps []Endpoint) *Client {
	c := New("", "", "")
	c.SetEndpoints(eps)
	return c
}

func (c *Client) SetEndpoints(eps []Endpoint) {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]Endpoint, 0, len(eps))
	for i, e := range eps {
		e.BaseURL = strings.TrimRight(strings.TrimSpace(e.BaseURL), "/")
		e.APIKey = strings.TrimSpace(e.APIKey)
		e.Model = strings.TrimSpace(e.Model)
		e.Provider = strings.TrimSpace(e.Provider)
		if e.BaseURL == "" {
			continue
		}
		if e.Provider == "" {
			e.Provider = detectProvider(e.BaseURL)
		}
		if e.APIKey == "" {
			switch e.Provider {
			case "opencode-zen":
				e.APIKey = "public"
			case "vercel-ai-gateway":
				e.APIKey = "oidc"
			case "pollinations":
				e.APIKey = "anonymous"
			default:
				continue
			}
		}
		if e.Model == "" {
			e.Model = "openai-fast"
		}
		if e.Name == "" {
			e.Name = fmt.Sprintf("ep%d", i+1)
		}
		out = append(out, e)
	}
	c.eps = out
}

func detectProvider(baseURL string) string {
	u := strings.ToLower(baseURL)
	if strings.Contains(u, "opencode.ai/zen") {
		return "opencode-zen"
	}
	return "openai"
}

func (c *Client) Configured() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.eps) > 0
}

// Status is a redacted snapshot for /ready.
func (c *Client) Status() map[string]any {
	if c == nil {
		return map[string]any{"configured": false}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	names := make([]string, 0, len(c.eps))
	for _, e := range c.eps {
		names = append(names, e.Name+":"+e.Model)
	}
	healthy := c.lastOK != "" && c.lastErr == ""
	checked := ""
	if !c.lastCheck.IsZero() {
		checked = c.lastCheck.Format(time.RFC3339)
	}
	return map[string]any{
		"configured": len(c.eps) > 0,
		"endpoints":  names,
		"last_ok":    c.lastOK,
		"last_error": c.lastErr,
		"checked_at": checked,
		"healthy":    healthy,
	}
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model       string    `json:"model"`
	Messages    []message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content          string `json:"content"`
			ReasoningContent string `json:"reasoning_content"`
		} `json:"message"`
	} `json:"choices"`
}

// CompleteTemp allows controlling creativity (router uses low temp; narrative uses higher).
func (c *Client) CompleteTemp(ctx context.Context, system, user string, temperature float64) (string, error) {
	return c.complete(ctx, system, user, temperature, 1200)
}

// CompleteBudget is CompleteTemp with a token cap (intent router, pings).
func (c *Client) CompleteBudget(ctx context.Context, system, user string, temperature float64, maxTokens int) (string, error) {
	if maxTokens < 32 {
		maxTokens = 256
	}
	return c.complete(ctx, system, user, temperature, maxTokens)
}

func (c *Client) complete(ctx context.Context, system, user string, temperature float64, maxTokens int) (string, error) {
	if c == nil || !c.Configured() {
		return "", fmt.Errorf("llm client not configured")
	}
	if temperature < 0 {
		temperature = 0
	}
	eps := c.activeEndpoints()
	if len(eps) == 0 {
		return "", fmt.Errorf("all llm endpoints in cooldown")
	}
	var last error
	for _, ep := range eps {
		text, err := c.call(ctx, ep, system, user, temperature, maxTokens)
		if err == nil && strings.TrimSpace(text) != "" {
			c.markOK(ep.Name)
			return strings.TrimSpace(text), nil
		}
		if err != nil {
			last = err
			c.markFail(ep.Name, err)
		} else {
			last = fmt.Errorf("%s: empty response", ep.Name)
			c.markFail(ep.Name, last)
		}
	}
	if last == nil {
		last = fmt.Errorf("llm failed")
	}
	return "", last
}

// Ping checks the first healthy endpoint (for /ready).
func (c *Client) Ping(ctx context.Context) error {
	if c == nil || !c.Configured() {
		return fmt.Errorf("llm client not configured")
	}
	_, err := c.CompleteBudget(ctx, "Reply with exactly: ok", "ping", 0, 16)
	return err
}

func (c *Client) activeEndpoints() []Endpoint {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	out := make([]Endpoint, 0, len(c.eps))
	for _, e := range c.eps {
		if until, ok := c.failUntil[e.Name]; ok && now.Before(until) {
			continue
		}
		out = append(out, e)
	}
	// If all cooling down, try primary anyway.
	if len(out) == 0 && len(c.eps) > 0 {
		return []Endpoint{c.eps[0]}
	}
	return out
}

func (c *Client) markOK(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.failUntil, name)
	c.lastOK = name
	c.lastErr = ""
	c.lastCheck = time.Now()
}

func (c *Client) markFail(name string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cool := 25 * time.Second
	msg := err.Error()
	if strings.Contains(msg, "FreeTierError") || strings.Contains(msg, "customer_verification_required") ||
		strings.Contains(msg, "not supported") {
		cool = 10 * time.Minute
		for _, e := range c.eps {
			if e.Name == name {
				for _, sib := range c.eps {
					if sib.Provider == e.Provider {
						c.failUntil[sib.Name] = time.Now().Add(cool)
					}
				}
				break
			}
		}
	} else {
		c.failUntil[name] = time.Now().Add(cool)
	}
	c.lastErr = fmt.Sprintf("%s: %v", name, err)
	c.lastCheck = time.Now()
}

func (c *Client) call(ctx context.Context, ep Endpoint, system, user string, temperature float64, maxTokens int) (string, error) {
	if maxTokens <= 0 {
		maxTokens = 1200
	}
	body, _ := json.Marshal(chatRequest{
		Model: ep.Model,
		Messages: []message{
			{Role: "system", Content: system},
			{Role: "user", Content: user},
		},
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Stream:      false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep.BaseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	switch ep.Provider {
	case "vercel-ai-gateway":
		key := strings.TrimSpace(os.Getenv("AI_GATEWAY_API_KEY"))
		if key == "" {
			key = strings.TrimSpace(os.Getenv("VERCEL_OIDC_TOKEN"))
		}
		if key == "" {
			key = ep.APIKey
		}
		if key == "" || key == "oidc" {
			return "", fmt.Errorf("vercel ai gateway: no token")
		}
		req.Header.Set("Authorization", "Bearer "+key)
		ua := c.UserAgent
		if ua == "" {
			ua = defaultUserAgent
		}
		req.Header.Set("User-Agent", ua)
	case "opencode-zen":
		// Free OpenCode Zen tier: public bearer + session headers (OpenAI-compatible).
		key := ep.APIKey
		if key == "" {
			key = "public"
		}
		req.Header.Set("Authorization", "Bearer "+key)
		req.Header.Set("User-Agent", openCodeZenUA)
		req.Header.Set("x-opencode-client", "cli")
		req.Header.Set("x-opencode-project", "global")
		req.Header.Set("x-opencode-request", "msg_"+randHex(8))
		req.Header.Set("x-opencode-session", "ses_"+randHex(8))
	case "pollinations":
		req.Header.Set("Authorization", "Bearer "+ep.APIKey)
		ua := c.UserAgent
		if ua == "" {
			ua = defaultUserAgent
		}
		req.Header.Set("User-Agent", ua)
		req.Header.Set("Referer", "https://pollinations.ai/")
	default:
		req.Header.Set("Authorization", "Bearer "+ep.APIKey)
		ua := c.UserAgent
		if ua == "" {
			ua = defaultUserAgent
		}
		req.Header.Set("User-Agent", ua)
	}
	res, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 2<<20))
	if ep.Provider == "pollinations" && (res.StatusCode == 402 || res.StatusCode == 403 || res.StatusCode == 429) {
		if text, gerr := c.pollinationsGet(ctx, system, user); gerr == nil && strings.TrimSpace(text) != "" {
			return strings.TrimSpace(text), nil
		} else if gerr != nil {
			return "", fmt.Errorf("status %d; get: %v", res.StatusCode, gerr)
		}
	}
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("status %d: %s", res.StatusCode, truncate(string(raw), 240))
	}
	var out chatResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("empty llm response")
	}
	msg := out.Choices[0].Message
	text := strings.TrimSpace(msg.Content)
	if text == "" {
		text = strings.TrimSpace(msg.ReasoningContent)
	}
	if text == "" {
		return "", fmt.Errorf("empty llm response")
	}
	return text, nil
}

// pollinationsGet is the anonymous text API. The OpenAI-compatible POST
// often answers 402 from Vercel while this route still returns text.
func (c *Client) pollinationsGet(ctx context.Context, system, user string) (string, error) {
	prompt := strings.TrimSpace(user)
	if i := strings.Index(prompt, "QUESTION:"); i >= 0 {
		prompt = strings.TrimSpace(prompt[i:])
	}
	prompt = clipRunes(prompt, 280)
	sys := clipRunes(strings.TrimSpace(system), 180)
	if prompt == "" {
		return "", fmt.Errorf("empty prompt")
	}
	u := "https://text.pollinations.ai/" + url.PathEscape(prompt) + "?model=openai-fast"
	if sys != "" {
		u += "&system=" + url.QueryEscape(sys)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	ua := c.UserAgent
	if ua == "" {
		ua = defaultUserAgent
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/plain")
	res, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 64<<10))
	if res.StatusCode >= 300 {
		return "", fmt.Errorf("pollinations get status %d: %s", res.StatusCode, truncate(string(raw), 180))
	}
	text := strings.TrimSpace(string(raw))
	if text == "" || strings.HasPrefix(text, "{") || strings.HasPrefix(text, "<") {
		return "", fmt.Errorf("pollinations get empty")
	}
	return text, nil
}

func clipRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}
