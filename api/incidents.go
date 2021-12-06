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
	IncidentStatus string
	Impact         string

	// All component statuses constants

	ReqIncident struct {
		Incident `json:"incident,omitempty"`
	}

	Incident struct {
		ID                                        string         `json:"id,omitempty"`
		Components                                interface{}    `json:"components,omitempty"`
		ComponentIDs                              []string       `json:"component_ids,omitempty"`
		CreatedAt                                 *time.Time     `json:"created_at,omitempty"`
		Impact                                    Impact         `json:"impact,omitempty"`
		ImpactOverride                            Impact         `json:"impact_override,omitempty"`
		IncidentUpdates                           []interface{}  `json:"incident_updates,omitempty"`
		Metadata                                  Metadata       `json:"metadata,omitempty"`
		MonitoringAt                              *time.Time     `json:"monitoring_at,omitempty"`
		Name                                      string         `json:"name,omitempty"`
		PageID                                    string         `json:"page_id,omitempty"`
		PostmortemBody                            string         `json:"postmortem_body,omitempty"`
		PostmortemBodyLastUpdatedAt               *time.Time     `json:"postmortem_body_last_updated_at,omitempty"`
		PostmortemIgnored                         bool           `json:"postmortem_ignored,omitempty"`
		PostmortemNotifiedSubscribers             bool           `json:"postmortem_notified_subscribers,omitempty"`
		PostmortemNotifiedTwitter                 bool           `json:"postmortem_notified_twitter,omitempty"`
		PostmortemPublishedAt                     bool           `json:"postmortem_published_at,omitempty"`
		ResolvedAt                                *time.Time     `json:"resolved_at,omitempty"`
		ScheduledAutoCompleted                    bool           `json:"scheduled_auto_completed,omitempty"`
		ScheduledAutoInProgress                   bool           `json:"scheduled_auto_in_progress,omitempty"`
		ScheduledAutoTransition                   bool           `json:"scheduled_auto_transition,omitempty"`
		ScheduledFor                              *time.Time     `json:"scheduled_for,omitempty"`
		ScheduledRemindPrior                      bool           `json:"scheduled_remind_prior,omitempty"`
		DeliverNotifications                      bool           `json:"deliver_notifications,omitempty"`
		AutoTransitionDeliverNotificationsAtEnd   bool           `json:"auto_transition_deliver_notifications_at_end,omitempty"`
		AutoTransitionDeliverNotificationsAtStart bool           `json:"auto_transition_deliver_notifications_at_start,omitempty"`
		AutoTransitionToMaintenanceState          bool           `json:"auto_transition_to_maintenance_state,omitempty"`
		AutoTransitionToOperationalState          bool           `json:"auto_transition_to_operational_state,omitempty"`
		AutoTweetAtBeginning                      bool           `json:"auto_tweet_at_beginning,omitempty"`
		AutoTweetOnCompletion                     bool           `json:"auto_tweet_on_completion,omitempty"`
		AutoTweetOnCreation                       bool           `json:"auto_tweet_on_creation,omitempty"`
		AutoTweetOneHourBefore                    bool           `json:"auto_tweet_one_hour_before,omitempty"`
		BackFillDate                              string         `json:"backfill_date,omitempty"`
		BackFilled                                bool           `json:"backfilled,omitempty"`
		Body                                      string         `json:"body,omitempty"`
		ScheduledRemindedAt                       *time.Time     `json:"scheduled_reminded_at,omitempty"`
		ScheduledUntil                            *time.Time     `json:"scheduled_until,omitempty"`
		Shortlink                                 string         `json:"shortlink,omitempty"`
		Status                                    IncidentStatus `json:"status,omitempty"`
		UpdatedAt                                 *time.Time     `json:"updated_at,omitempty"`
	}

	Metadata struct {
	}
)

const (
	IncidentStatusInvestigating IncidentStatus = "investigating"
	IncidentStatusIdentified    IncidentStatus = "identified"
	IncidentStatusMonitoring    IncidentStatus = "monitoring"
	IncidentStatusResolved      IncidentStatus = "resolved"
	IncidentStatusScheduled     IncidentStatus = "scheduled"
	IncidentStatusInProgress    IncidentStatus = "in_progress"
	IncidentStatusVerifying     IncidentStatus = "verifying"
	IncidentStatusCompleted     IncidentStatus = "completed"

	ImpactMaintenance Impact = "maintenance"
	ImpactCritical    Impact = "critical"
	ImpactNone        Impact = "none"
	ImpactMajor       Impact = "major"
	ImpactMinor       Impact = "minor"
)

func (i IncidentStatus) String() string {

	return string(i)
}

func (i Impact) String() string {

	return string(i)
}

func (s StatusPage) GetIncidents() ([]Incident, error) {

	var incidents []Incident
	url := fmt.Sprintf("%s/v1/pages/%s/incidents", s.Client.Config.URL, s.Page.ID)

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error %s", err)
		return incidents, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return incidents, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return incidents, err
	}
	if rsp.StatusCode != 200 {
		return incidents, fmt.Errorf("error %s %s", rsp.Status, body)
	}

	json.Unmarshal(body, &incidents)

	return incidents, nil

}

func (s StatusPage) GetUnresolvedIncidents() ([]Incident, error) {

	var incidents []Incident
	url := fmt.Sprintf("%s/v1/pages/%s/incidents/unresolved", s.Client.Config.URL, s.Page.ID)

	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error %s", err)
		return incidents, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return incidents, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return incidents, err
	}
	if rsp.StatusCode != 200 {
		return incidents, fmt.Errorf("error %s %s", rsp.Status, body)
	}

	json.Unmarshal(body, &incidents)

	return incidents, nil

}

func (s StatusPage) UpdateIncident(i Incident) (Incident, error) {
	incident := Incident{
		Name:                                    i.Name,
		Status:                                  i.Status,
		ScheduledFor:                            i.ScheduledFor,
		ScheduledUntil:                          i.ScheduledUntil,
		ScheduledRemindPrior:                    i.ScheduledRemindPrior,
		ScheduledAutoInProgress:                 i.ScheduledAutoInProgress,
		ScheduledAutoCompleted:                  i.ScheduledAutoCompleted,
		Metadata:                                i.Metadata,
		DeliverNotifications:                    i.DeliverNotifications,
		AutoTransitionDeliverNotificationsAtEnd: i.AutoTransitionDeliverNotificationsAtEnd,
		AutoTransitionDeliverNotificationsAtStart: i.AutoTransitionDeliverNotificationsAtStart,
		AutoTransitionToMaintenanceState:          i.AutoTransitionToMaintenanceState,
		AutoTransitionToOperationalState:          i.AutoTransitionToOperationalState,
		AutoTweetAtBeginning:                      i.AutoTweetAtBeginning,
		AutoTweetOnCompletion:                     i.AutoTweetOnCompletion,
		AutoTweetOnCreation:                       i.AutoTweetOnCreation,
		AutoTweetOneHourBefore:                    i.AutoTweetOneHourBefore,
		BackFillDate:                              i.BackFillDate,
		BackFilled:                                i.BackFilled,
		Body:                                      i.Body,
		Components:                                i.Components,
		ComponentIDs:                              i.ComponentIDs,
	}

	url := fmt.Sprintf("%s/v1/pages/%s/incidents/%s", s.Client.Config.URL, s.Page.ID, i.ID)

	recComp := ReqIncident{Incident: incident}
	b, err := json.Marshal(recComp)
	if err != nil {
		return i, err
	}
	r, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Error %s", err)
		return i, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	r.Header.Add("Content-Type", "application/json")
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return i, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return i, err
	}
	if rsp.StatusCode != 200 {
		return i, fmt.Errorf("error %s %s", rsp.Status, body)
	}

	json.Unmarshal(body, &i)

	return i, nil

}

func (s StatusPage) CreateIncident(i Incident) (Incident, error) {
	incident := Incident{
		Name:           i.Name,
		Status:         i.Status,
		ImpactOverride: i.ImpactOverride,
		Body:           i.Body,
		Components:     i.Components,
		ComponentIDs:   i.ComponentIDs,
	}
	url := fmt.Sprintf("%s/v1/pages/%s/incidents", s.Client.Config.URL, s.Page.ID)

	recComp := ReqIncident{Incident: incident}
	b, err := json.Marshal(recComp)
	if err != nil {
		return i, err
	}
	fmt.Println(string(b))
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Error %s", err)
		return i, err
	}
	r.Header.Add("Authorization", fmt.Sprintf("OAuth %s", s.Client.Config.Token))
	r.Header.Add("Content-Type", "application/json")
	rsp, err := s.Client.httpclient.Do(r)
	if err != nil {
		log.Printf("Error %s", err)
		return i, err
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Printf("Error %s", err)
		return i, err
	}
	if rsp.StatusCode != 201 {
		return i, fmt.Errorf("error %s %s", rsp.Status, body)
	}

	json.Unmarshal(body, &i)

	return i, nil

}

func (s StatusPage) DeleteIncident(i Incident) error {

	url := fmt.Sprintf("%s/v1/pages/%s/incidents/%s", s.Client.Config.URL, s.Page.ID, i.ID)

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

func (s StatusPage) FilterIncidents(componentID string, status IncidentStatus) ([]Incident, error) {

	var incidents []Incident

	if status != IncidentStatusResolved {
		i, err := s.GetUnresolvedIncidents()
		if err != nil {
			return incidents, fmt.Errorf("unable to get incidents %s", err)
		}

		for _, incident := range i {

			if incident.Status == status {

				for _, id := range incident.ComponentIDs {

					if id == componentID {

						incidents = append(incidents, incident)

					}

				}

			}
		}
	} else {

		i, err := s.GetIncidents()
		if err != nil {
			return incidents, fmt.Errorf("unable to get incidents %s", err)
		}

		for _, incident := range i {

			if incident.Status == status {

				for _, id := range incident.ComponentIDs {

					if id == componentID {

						incidents = append(incidents, incident)

					}

				}

			}
		}

	}

	return incidents, nil

}

func (s StatusPage) GetOpenedIncidentByName(name, componentID string) (i Incident, err error) {
	var incidents []Incident

	incidents, err = s.GetUnresolvedIncidents()
	if err != nil {
		return i, fmt.Errorf("unable to get incidents %s", err)
	}

	for _, incident := range incidents {

		if incident.Name == name {

			for _, id := range incident.ComponentIDs {

				if id == componentID {

					return incident, nil

				}

			}

		}

	}
	return i, fmt.Errorf("unable to find incident by name", name)

}
