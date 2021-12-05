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
	ReqComponentGroup struct {
		Description    string `json:"description,omitempty"`
		ComponentGroup `json:"component_group,omitempty"`
	}

	ComponentGroup struct {
		ID          string     `json:"id,omitempty"`
		PageID      string     `json:"page_id,omitempty"`
		Name        string     `json:"name,omitempty"`
		Description string     `json:"description,omitempty"`
		Components  []string   `json:"components"`
		Position    string     `json:"position,omitempty"`
		CreatedAt   *time.Time `json:"created_at,omitempty"`
		UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	}
)

func (s StatusPage) GetComponentGroups() ([]ComponentGroup, error) {

	groups := []ComponentGroup{}
	url := fmt.Sprintf("%s/v1/pages/%s/component-groups", s.Client.Config.URL, s.Page.ID)

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error %s", err)
		return groups, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return groups, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return groups, err
	}
	if rsp.StatusCode != 200 {
		return groups, fmt.Errorf("error %s %s", rsp.Status, body)
	}
	json.Unmarshal(body, &groups)

	return groups, nil

}

func (s StatusPage) UpdateComponentGroup(c ComponentGroup) (ComponentGroup, error) {

	compGroup := ComponentGroup{
		Components: c.Components,
		Name:       c.Name,
	}
	url := fmt.Sprintf("%s/v1/pages/%s/component-groups/%s", s.Client.Config.URL, s.Page.ID, c.ID)
	cg := ReqComponentGroup{Description: c.Description, ComponentGroup: compGroup}

	b, err := json.Marshal(&cg)
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

func (s StatusPage) CreateComponentGroup(c ComponentGroup) (ComponentGroup, error) {

	url := fmt.Sprintf("%s/v1/pages/%s/component-groups", s.Client.Config.URL, s.Page.ID)
	cg := ReqComponentGroup{Description: c.Description, ComponentGroup: c}

	b, err := json.Marshal(&cg)
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

func (s StatusPage) DeleteComponentGroups(c ComponentGroup) error {

	url := fmt.Sprintf("%s/v1/pages/%s/component-groups/%s", s.Client.Config.URL, s.Page.ID, c.ID)

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

func (s StatusPage) GetComponentGroupByName(name string) (c ComponentGroup, err error) {
	groups, err := s.GetComponentGroups()
	if err != nil {
		log.Printf("Error %s", err)
		return c, err

	}

	for _, g := range groups {
		if g.Name == name {

			return g, nil

		}

	}

	return c, fmt.Errorf("unable to find group %s", name)

}
