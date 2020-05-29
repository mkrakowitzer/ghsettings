package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mkrakowitzer/ghsettings/config"
)

type Collaborator struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

func UpdateCollaborator(apiClient *Client, config config.C) error {

	err := CollaboratorAddToRepo(apiClient, config)
	if err != nil {
		log.Fatal(err)
	}

	err = CollaboratorRemoveFromRepo(apiClient, config)
	return err
}

func CollaboratorAddToRepo(client *Client, config config.C) error {

	for _, s := range config.Collaborators {

		path := fmt.Sprintf("repos/%s/%s/collaborators/%s", Org, config.Repository.Name, s.Username)
		result := Team{}

		collaborator := Collaborator{
			Permission: s.Permission,
		}

		j, _ := json.Marshal(collaborator)

		err := client.REST("PUT", path, bytes.NewBuffer(j), &result)
		if err != nil {
			return err
		}
	}
	return nil
}

type ListCollaborator struct {
	Affiliation string `json:"affiliation"`
}

type Collaborators []struct {
	Login  string `json:"login"`
	ID     int    `json:"id"`
	NodeID string `json:"node_id"`
	Type   string `json:"type"`
}

// Refactor this, was in a hurry
func CollaboratorRemoveFromRepo(client *Client, config config.C) error {

	path := fmt.Sprintf("repos/%s/%s/collaborators", Org, config.Repository.Name)
	result := Collaborators{}

	err := client.REST("GET", path, &bytes.Buffer{}, &result)
	if err != nil {
		return err
	}

	var gh_rules []string
	var yml_rules []string
	for _, k := range config.Collaborators {
		yml_rules = append(yml_rules, k.Username)
	}
	for _, k := range result {
		if k.Type == "User" {
			gh_rules = append(gh_rules, k.Login)
		}
	}

	delete := missing(yml_rules, gh_rules)

	for _, k := range result {
		for _, s := range delete {
			if k.Login == s {
				result1 := Teams{}
				path := fmt.Sprintf("repos/%s/%s/collaborators/%s", Org, config.Repository.Name, k.Login)
				err := client.REST("DELETE", path, &bytes.Buffer{}, &result1)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
