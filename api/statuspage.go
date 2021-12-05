package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	StatusPage struct {
		Client *Client
		Page
	}

	Pages []struct {
		Page
	}

	Page struct {
		ID                       string    `json:"id"`
		CreatedAt                time.Time `json:"created_at"`
		UpdatedAt                time.Time `json:"updated_at"`
		Name                     string    `json:"name"`
		PageDescription          string    `json:"page_description"`
		Headline                 string    `json:"headline"`
		Branding                 string    `json:"branding"`
		Subdomain                string    `json:"subdomain"`
		Domain                   string    `json:"domain"`
		URL                      string    `json:"url"`
		SupportURL               string    `json:"support_url"`
		HiddenFromSearch         bool      `json:"hidden_from_search"`
		AllowPageSubscribers     bool      `json:"allow_page_subscribers"`
		AllowIncidentSubscribers bool      `json:"allow_incident_subscribers"`
		AllowEmailSubscribers    bool      `json:"allow_email_subscribers"`
		AllowSmsSubscribers      bool      `json:"allow_sms_subscribers"`
		AllowRssAtomFeeds        bool      `json:"allow_rss_atom_feeds"`
		AllowWebhookSubscribers  bool      `json:"allow_webhook_subscribers"`
		NotificationsFromEmail   string    `json:"notifications_from_email"`
		NotificationsEmailFooter string    `json:"notifications_email_footer"`
		ActivityScore            int       `json:"activity_score"`
		TwitterUsername          string    `json:"twitter_username"`
		ViewersMustBeTeamMembers bool      `json:"viewers_must_be_team_members"`
		IPRestrictions           string    `json:"ip_restrictions"`
		City                     string    `json:"city"`
		State                    string    `json:"state"`
		Country                  string    `json:"country"`
		TimeZone                 string    `json:"time_zone"`
		CSSBodyBackgroundColor   string    `json:"css_body_background_color"`
		CSSFontColor             string    `json:"css_font_color"`
		CSSLightFontColor        string    `json:"css_light_font_color"`
		CSSGreens                string    `json:"css_greens"`
		CSSYellows               string    `json:"css_yellows"`
		CSSOranges               string    `json:"css_oranges"`
		CSSBlues                 string    `json:"css_blues"`
		CSSReds                  string    `json:"css_reds"`
		CSSBorderColor           string    `json:"css_border_color"`
		CSSGraphColor            string    `json:"css_graph_color"`
		CSSLinkColor             string    `json:"css_link_color"`
		CSSNoData                string    `json:"css_no_data"`
		FaviconLogo              string    `json:"favicon_logo"`
		TransactionalLogo        string    `json:"transactional_logo"`
		HeroCover                string    `json:"hero_cover"`
		EmailLogo                string    `json:"email_logo"`
		TwitterLogo              string    `json:"twitter_logo"`
	}
	Client struct {
		Config     *Config
		httpclient *http.Client
	}

	Config struct {
		URL     string
		Token   string
		Timeout time.Duration
	}
)

func New(url, token string, timeout time.Duration) StatusPage {
	c := &Config{URL: url,
		Token:   token,
		Timeout: timeout,
	}

	return StatusPage{
		Client: &Client{
			Config:     c,
			httpclient: &http.Client{},
		},
		Page: Page{},
	}
}

func (c Client) GetPages() (Pages, error) {

	pages := Pages{}
	url := fmt.Sprintf("%s/v1/pages", c.Config.URL)

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error %s", err)
		return pages, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", c.Config.Token))
	rsp, err := c.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return pages, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return pages, err
	}

	json.Unmarshal(body, &pages)

	return pages, nil

}
