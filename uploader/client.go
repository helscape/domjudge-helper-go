package uploader

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Client struct {
	serverURL     string
	username      string
	password      string
	actionTimeout time.Duration

	browserOnce      sync.Once
	browser          *rod.Browser
	activePage       *rod.Page
	browserMu        sync.Mutex
	browserCloseOnce sync.Once

	probNameCache map[string]string
	probNameMu    sync.RWMutex
}

func NewClientWithTimeout(serverURL, username, password string, timeoutSeconds int) *Client {
	d := time.Duration(timeoutSeconds) * time.Second
	return &Client{
		serverURL:     serverURL,
		username:      username,
		password:      password,
		actionTimeout: d,
		probNameCache: make(map[string]string),
	}
}

func (c *Client) initBrowser() error {
	var initErr error
	c.browserOnce.Do(func() {
		u, err := launcher.New().
			Headless(true).
			Leakless(false).
			Launch()
		if err != nil {
			initErr = fmt.Errorf("cannot launch browser: %w", err)
			return
		}

		c.browser = rod.New().ControlURL(u).MustConnect()

		page, err := c.browser.Page(proto.TargetCreateTarget{URL: "about:blank"})
		if err != nil {
			initErr = fmt.Errorf("open page: %w", err)
			return
		}

		if err := c.browserLogin(page); err != nil {
			initErr = fmt.Errorf("login: %w", err)
			return
		}

		c.activePage = page
	})
	return initErr
}

func (c *Client) CloseBrowser() {
	c.browserCloseOnce.Do(func() {
		if c.activePage != nil {
			c.activePage.Close()
		}
		if c.browser != nil {
			c.browser.MustClose()
		}
	})
}

func (c *Client) GetProblemName(probID string) (string, error) {
	c.probNameMu.RLock()
	if name, ok := c.probNameCache[probID]; ok {
		c.probNameMu.RUnlock()
		return name, nil
	}
	c.probNameMu.RUnlock()

	if err := c.initBrowser(); err != nil {
		return "", err
	}

	c.browserMu.Lock()
	defer c.browserMu.Unlock()

	p := c.activePage.Timeout(c.actionTimeout)
	url := fmt.Sprintf("%s/jury/problems/%s/testcases", c.serverURL, probID)

	if err := p.Navigate(url); err != nil {
		return "", fmt.Errorf("navigate to scrape: %w", err)
	}

	el, err := p.Element("h1")
	if err != nil {
		return "", fmt.Errorf("find h1 element: %w", err)
	}

	rawText, err := el.Text()
	if err != nil {
		return "", fmt.Errorf("extract h1 text: %w", err)
	}

	name := rawText
	if idx := strings.Index(rawText, " - "); idx != -1 {
		name = strings.TrimSpace(rawText[idx+3:])
	}

	c.probNameMu.Lock()
	c.probNameCache[probID] = name
	c.probNameMu.Unlock()

	return name, nil
}

func (c *Client) UploadTestcaseBrowser(probID, inputFile, outputFile string) error {
	if err := c.initBrowser(); err != nil {
		return err
	}

	c.browserMu.Lock()
	defer c.browserMu.Unlock()

	p := c.activePage.Timeout(c.actionTimeout)
	uploadURL := fmt.Sprintf("%s/jury/problems/%s/testcases", c.serverURL, probID)
	return c.browserUpload(p, uploadURL, inputFile, outputFile)
}
