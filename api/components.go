package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	ComponentStatus string
	ReqComponent    struct {
		Component Component `json:"component,omitempty"`
	}

	Component struct {
		ID                 string          `json:"id,omitempty"`
		PageID             string          `json:"page_id,omitempty"`
		GroupID            string          `json:"group_id,omitempty"`
		CreatedAt          *time.Time      `json:"created_at,omitempty"`
		UpdatedAt          *time.Time      `json:"updated_at,omitempty"`
		Group              bool            `json:"group,omitempty"`
		Name               string          `json:"name,omitempty"`
		Description        string          `json:"description,omitempty"`
		Position           int             `json:"position,omitempty"`
		Status             ComponentStatus `json:"status,omitempty"`
		Showcase           bool            `json:"showcase,omitempty"`
		OnlyShowIfDegraded bool            `json:"only_show_if_degraded,omitempty"`
		AutomationEmail    string          `json:"automation_email,omitempty"`
		StartDate          string          `json:"start_date,omitempty"`
	}
)

const (
	ComponentStatusOperational         ComponentStatus = "operational"
	ComponentStatusUnderMaintenance    ComponentStatus = "under_maintenance"
	ComponentStatusDegradedPerformance ComponentStatus = "degraded_performance"
	ComponentStatusPartialOutage       ComponentStatus = "partial_outage"
	ComponentStatusMajorOutage         ComponentStatus = "major_outage"
	ComponentStatusEmpty               ComponentStatus = ""
)

func (c ComponentStatus) String() string {

	return string(c)
}
func (s StatusPage) GetComponents() ([]Component, error) {

	var components []Component
	url := fmt.Sprintf("%s/v1/pages/%s/components", s.Client.Config.URL, s.Page.ID)

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error %s", err)
		return components, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return components, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return components, err
	}
	if rsp.StatusCode != 200 {
		return components, fmt.Errorf("error %s %s", rsp.Status, body)
	}

	json.Unmarshal(body, &components)

	return components, nil

}

func (s StatusPage) UpdateComponent(c Component) (Component, error) {
	comp := Component{
		Description:        c.Description,
		Status:             c.Status,
		Name:               c.Name,
		OnlyShowIfDegraded: c.OnlyShowIfDegraded,
		GroupID:            c.GroupID,
		Showcase:           c.Showcase,
		StartDate:          c.StartDate,
	}
	url := fmt.Sprintf("%s/v1/pages/%s/components/%s", s.Client.Config.URL, s.Page.ID, c.ID)

	recComp := ReqComponent{Component: comp}
	b, err := json.Marshal(recComp)
	if err != nil {
		return c, err
	}
	r, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Error %s", err)
		return c, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	r.Header.Add("Content-Type", "application/json")
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return c, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return c, err
	}
	if rsp.StatusCode != 200 {
		return c, fmt.Errorf("error %s %s", rsp.Status, body)
	}

	json.Unmarshal(body, &c)

	return c, nil

}

func (s StatusPage) CreateComponent(c Component) (Component, error) {

	url := fmt.Sprintf("%s/v1/pages/%s/components", s.Client.Config.URL, s.Page.ID)

	recComp := ReqComponent{Component: c}
	b, err := json.Marshal(recComp)
	log.Print(string(b))
	log.Print(url)
	if err != nil {
		return c, err
	}
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Error %s", err)
		return c, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	r.Header.Add("Content-Type", "application/json")
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return c, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return c, err
	}
	if rsp.StatusCode != 201 {
		return c, fmt.Errorf("error %s %s", rsp.Status, body)
	}

	json.Unmarshal(body, &c)

	return c, nil

}

func (s StatusPage) DeleteComponent(c Component) error {

	url := fmt.Sprintf("%s/v1/pages/%s/components/%s", s.Client.Config.URL, s.Page.ID, c.ID)

	r, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return err
	}
	if rsp.StatusCode != 200 {
		return fmt.Errorf("error %s %s", rsp.Status, body)
	}

	return nil

}

// GetComponentID will check if given Component already exists on gid
func (s StatusPage) GetComponentByName(name string, gid string) (c Component, err error) {
	components, err := s.GetComponents()
	if err != nil {
		log.Printf("Error %s", err)
		return c, err
	}
	for _, comp := range components {

		if comp.Name == name && comp.GroupID == gid {

			return comp, nil

		}

	}

	return c, fmt.Errorf("unable find component %s", name)

}
