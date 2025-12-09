package airbnb

import (
	"fmt"

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
	c := &config{
		RodURL:   "",
		HeadLess: false,
	}
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
		browser:                            browser,
		hasGonePastTheTheTranslationDialog: false,
	}, nil
}

// Close closes the browser connection.
func (c *Client) Close() error {
	if c.browser != nil {
		err := c.browser.Close()
		if err != nil {
			return fmt.Errorf("failed to close browser: %w", err)
		}
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
			return nil, fmt.Errorf("failed to connect to browser: %w", err)
		}
		return browser, nil
	}

	// Connect to remote Rod service
	browser := rod.New().ControlURL("ws://" + config.RodURL)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to remote Rod service: %w", err)
	}
	return browser, nil
}
