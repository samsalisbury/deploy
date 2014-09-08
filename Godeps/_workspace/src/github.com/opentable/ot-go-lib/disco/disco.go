// Package disco is the client for disco servers like http://github.com/opentable/disco
package disco

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opentable/ot-go-lib/env"
	"github.com/opentable/ot-go-lib/logging"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	discoURL = env.RequireString("OT_CLOUD_PLATFORM_DISCO_URL")
	host     = env.RequireString("APP_HOST")
	port     = env.RequireString("PORT0")
)

type Client struct {
	updateReq     *http.Request
	announceURL   string
	unannounceReq *http.Request
	comment       string
	htc           http.Client
	reg           Registry
	log           logging.StartupLog
	stop          chan struct{}
}

var serviceTypePattern = regexp.MustCompile(`[a-z0-9\-]+`)

// NewClientFromEnv creates a client that knows how to announce and unannounce the service
// named serviceType at the URL appURL.
// It expects an env var OT_CLOUD_PLATFORM_DISCO_URL in order to work.
func NewClientFromEnv(serviceType string, comment string) (*Client, error) {
	appURL := "http://" + host + ":" + port
	if strings.HasPrefix(appURL, "http://") {
		appURL = "http:" + strings.TrimPrefix(appURL, "http://")
	} else if strings.HasPrefix(appURL, "https://") {
		appURL = "https:" + strings.TrimPrefix(appURL, "https://")
	} else {
		return nil, fmt.Errorf("appURL must begin either http:// or https://")
	}
	appURL = strings.TrimRight(appURL, "/")
	if !serviceTypePattern.MatchString(serviceType) {
		return nil, fmt.Errorf(`serviceType should match the expression [a-z0-9\-]+`)
	}
	announceURL := discoURL + "/" + serviceType + "/" + appURL

	updateReq, err := http.NewRequest("GET", discoURL, nil)
	if err != nil {
		return nil, err
	}
	updateReq.Header.Add("Accept", "application/json")

	unannounceReq, err := http.NewRequest("DELETE", announceURL, nil)

	log := logging.StandardConfig(serviceType + "-discoclient").StartupLog(0)
	stop := make(chan struct{})
	return &Client{updateReq, announceURL, unannounceReq, comment, http.Client{}, Registry{}, log, stop}, nil
}

func (c *Client) Announce() {
	req, err := http.NewRequest("PUT", c.announceURL, bytes.NewReader([]byte(c.comment)))
	if err != nil {
		c.log.Error("Announce failed to create HTTP request:", err)
		return
	}
	req.Header.Add("Content-Type", "text/plain")
	r, err := c.htc.Do(req)
	if err != nil {
		c.log.Error("Announce failed:", err)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	if r.StatusCode != 200 && r.StatusCode != 201 {
		c.log.Error("Announce got status code", r.StatusCode, "want 200 or 201 ("+req.URL.String()+")")
	}
}

func (c *Client) AnnounceEvery500ms() {
	for {
		select {
		case <-time.After(500 * time.Millisecond):
			c.Announce()
			err := c.Update()
			if err != nil {
				c.log.Error(err)
			}
		case <-c.stop:
			return
		}
	}
}

func (c *Client) Unannounce() {
	c.stop <- struct{}{}
	r, err := c.htc.Do(c.unannounceReq)
	if err != nil {
		c.log.Error("Unannounce failed:", err)
		return
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	if r.StatusCode != 204 {
		c.log.Error("Unannounce got status code", r.StatusCode, "want 204 ("+c.unannounceReq.URL.String()+")")
	}
}

func (c *Client) Update() error {
	r, err := c.htc.Do(c.updateReq)
	if err != nil {
		return fmt.Errorf("Update failed: %v", err)
	}
	if r.Body != nil {
		defer r.Body.Close()
	}
	if r.StatusCode != 200 {
		return fmt.Errorf("Update got status code: %v, want 200.", r.StatusCode)
	}
	var reg Registry
	err = json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {
		return fmt.Errorf("Update unable to deserialise response: %v", err)
	}
	c.reg = reg
	return nil
}

func (c *Client) Discover(serviceType string) (string, error) {
	a, ok := c.reg[serviceType]
	if !ok || len(a) == 0 {
		return "", fmt.Errorf("service '%v' not found", serviceType)
	}
	// Go always randomises map iterations... So we just need to grap the first one it picks
	for _, s := range a {
		return s.URL, nil
	}
	panic("You have reached unreachable code.")
}

type Registry map[string]map[string]*Service

type Service struct {
	SelfURL     string    `json:"_self"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Comment     string    `json:"comment"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	Expires     time.Time `json:"expires"`
	UpdateCount uint64    `json:"updateCount"`
}
