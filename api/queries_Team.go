package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mkrakowitzer/ghsettings/config"
	"github.com/mkrakowitzer/ghsettings/utils"
)

var Org string

type Team struct {
	Name       string `json:"name"`
	Permission string `json:"permission"`
}

func UpdateTeam(apiClient *Client, config config.C) error {
	err := TeamAddToRepo(apiClient, config)
	if err != nil {
		log.Fatal(err)
	}
	err = TeamDeleteFromRepo(apiClient, config)
	return err
}

func TeamAddToRepo(client *Client, config config.C) error {

	for _, s := range config.Teams {

		path := fmt.Sprintf("orgs/%s/teams/%s/repos/%s/%s", Org, s.Name, Org, config.Repository.Name)
		result := Team{}

		team := Team{
			Permission: s.Permission,
		}

		j, _ := json.Marshal(team)

		err := client.REST("PUT", path, bytes.NewBuffer(j), &result)
		if err != nil {
			return err
		}
	}
	return nil
}

type Teams []struct {
	ID     int    `json:"id"`
	NodeID string `json:"node_id"`
	Name   string `json:"name"`
	Slug   string `json:"slug"`
}

// Refactor this, was in a hurry
func TeamDeleteFromRepo(client *Client, config config.C) error {

	path := fmt.Sprintf("repos/%s/%s/teams", Org, config.Repository.Name)
	result := Teams{}

	err := client.REST("GET", path, &bytes.Buffer{}, &result)
	if err != nil {
		return err
	}

	var gh_rules []string
	var yml_rules []string
	for _, k := range config.Teams {
		yml_rules = append(yml_rules, k.Name)
	}
	for _, k := range result {
		gh_rules = append(gh_rules, k.Name)
	}

	delete := utils.Missing(yml_rules, gh_rules)

	for _, k := range result {
		for _, s := range delete {
			if k.Name == s {
				result1 := Teams{}
				path := fmt.Sprintf("orgs/%s/teams/%s/repos/%s/%s", Org, k.Slug, Org, config.Repository.Name)
				err := client.REST("DELETE", path, &bytes.Buffer{}, &result1)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
