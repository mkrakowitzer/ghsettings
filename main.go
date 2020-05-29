package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mkrakowitzer/ghsettings/api"
	"github.com/mkrakowitzer/ghsettings/config"
	"github.com/mkrakowitzer/ghsettings/context"
	"gopkg.in/yaml.v2"
)

var apiClientForContext = func(ctx context.Context) (*api.Client, error) {
	token, err := ctx.AuthToken()
	if err != nil {
		return nil, err
	}

	var opts []api.ClientOption

	getAuthValue := func() string {
		return fmt.Sprintf("token %s", token)
	}

	Version := "1"
	opts = append(opts,
		api.AddHeaderFunc("Authorization", getAuthValue),
		api.AddHeader("User-Agent", fmt.Sprintf("ghadmin %s", Version)),
		api.AddHeader("Accept", "application/vnd.github.antiope-preview+json"),
	)

	return api.NewClient(opts...), nil

}

var Org = os.Getenv("GITHUB_ORG")

func main() {
	fmt.Println("Running")
	api.Org = Org

	var config config.C

	files, err := ioutil.ReadDir("./repo_config")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.New()

	apiClient, err := apiClientForContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		data, err := ioutil.ReadFile("./repo_config/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Fatal(err)
		}

		repo, err := api.GetRepoID(apiClient, Org, config.Repository.Name)
		if err != nil {
			log.Fatal(err)
		}

		err = api.UpdateRepository(apiClient, repo, config)
		if err != nil {
			log.Fatal(err)
		}

		err = api.UpdateCollaborator(apiClient, config)
		if err != nil {
			log.Fatal(err)
		}

		err = api.UpdateTeam(apiClient, config)
		if err != nil {
			log.Fatal(err)
		}

		err = api.BranchProtections(apiClient, repo, config)
		if err != nil {
			log.Fatal(err)
		}
	}
}
