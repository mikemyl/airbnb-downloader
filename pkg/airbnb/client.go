package airbnb

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

// Client represents an Airbnb scraper client.
type Client struct {
	browser                            *rod.Browser
	hasGonePastTheTheTranslationDialog bool
}

type config struct {
	RodURL   string
	HeadLess bool
}

type Option func(*config) error

func WithHeadless(headLess bool) Option {
	return func(c *config) error {
		c.HeadLess = headLess
		return nil
	}
}

func WithRodURL(rodURL string) Option {
	return func(c *config) error {
		c.RodURL = rodURL
		return nil
	}
}

// NewClient creates a new Airbnb client with the given configuration.
func NewClient(opts ...Option) (*Client, error) {
	c := &config{}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	browser, err := createBrowser(c)
	if err != nil {
		return nil, err
	}

	return &Client{
		browser: browser,
	}, nil
}

// Close closes the browser connection.
func (c *Client) Close() error {
	if c.browser != nil {
		return c.browser.Close()
	}
	return nil
}

// createBrowser creates a browser instance based on the configuration.
func createBrowser(config *config) (*rod.Browser, error) {
	if config.RodURL == "" {
		u := launcher.New().Headless(config.HeadLess).MustLaunch()
		browser := rod.New().ControlURL(u)
		err := browser.Connect()
		if err != nil {
			return nil, err
		}
		return browser, nil
	}

	// Connect to remote Rod service
	browser := rod.New().ControlURL("ws://" + config.RodURL)
	if err := browser.Connect(); err != nil {
		return nil, err
	}
	return browser, nil
}
