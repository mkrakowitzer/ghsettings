package command

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/mkrakowitzer/ghsettings/api"
	"github.com/mkrakowitzer/ghsettings/config"
	"github.com/mkrakowitzer/ghsettings/context"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghsettings",
	Short: "Configure GitHub repositories, collaborators, teams and branch protections",
	Long: `Configure GitHub repositories, collaborators, teams and branch protections

MU_GITHUB_TOKEN and GITHUB_ORG environment variables must be set`,
	RunE: run,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ghsettings.yaml)")
	rootCmd.PersistentFlags().StringSlice("files", []string{}, "List of files seperated by spaces")
	viper.BindPFlag("files", rootCmd.PersistentFlags().Lookup("files"))
	viper.SetDefault("files", []string{})
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("enforce", "e", false, "Enforce Collaborators, Teams and Branches")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".ghsettings")
	}

	viper.BindEnv("GITHUB_ORG")
	viper.BindEnv("MU_GITHUB_TOKEN")
	viper.BindEnv("GHSETTINGS_CONFIGDIR")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run(cmd *cobra.Command, args []string) error {

	var config config.C

	Org := viper.GetString("GITHUB_ORG")
	api.Org = Org
	config_dir := viper.GetString("GHSETTINGS_CONFIGDIR")
	if config_dir == "" {
		config_dir = "repo_config"
	}

	var files []string

	if len(viper.GetStringSlice("files")) == 0 {
		f, err := ioutil.ReadDir(fmt.Sprintf("./%s", config_dir))
		if err != nil {
			log.Fatal(err)
		}
		for _, k := range f {
			files = append(files, fmt.Sprintf("%s/%s", config_dir, k.Name()))
		}
	} else {
		files = viper.GetStringSlice("files")
	}

	ctx := context.New()

	apiClient, err := apiClientForContext(ctx)
	if err != nil {
		log.Fatal(err)
	}
	rate_start, _ := api.GetRateLimit(apiClient)

	for _, f := range files {

		data, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Fatal(err)
		}
		log.WithFields(log.Fields{
			"name": config.Repository.Name,
		}).Info("applying to repository")

		repo, err := api.GetRepoID(apiClient, Org, config.Repository.Name)
		if err != nil {
			log.Fatal(err)
		}

		err = api.UpdateRepository(apiClient, repo, config)
		if err != nil {
			log.Fatal(err)
		}

		err = api.UpdateCollaborator(apiClient, config, cmd)
		if err != nil {
			log.Fatal(err)
		}

		err = api.UpdateTeam(apiClient, config, cmd)
		if err != nil {
			log.Fatal(err)
		}

		err = api.BranchProtections(apiClient, repo, config, cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
	rate_end, _ := api.GetRateLimit(apiClient)
	log.WithFields(log.Fields{
		"core_api_calls":     rate_start.Resources.Core.Remaining - rate_end.Resources.Core.Remaining,
		"graphql_ap_calls":   rate_start.Resources.Graphql.Remaining - rate_end.Resources.Graphql.Remaining,
		"core_remaining":     rate_end.Resources.Core.Remaining,
		"graphql_remaining":  rate_end.Resources.Graphql.Remaining,
		"combined_remaining": rate_end.Rate.Remaining,
	}).Info("rate limit stats")
	return nil
}

var apiClientForContext = func(ctx context.Context) (*api.Client, error) {
	token, err := ctx.AuthToken()
	if err != nil {
		return nil, err
	}

	var opts []api.ClientOption
	if verbose := os.Getenv("DEBUG"); verbose != "" {
		opts = append(opts, api.ApiVerboseLog())
	}
	getAuthValue := func() string {
		return fmt.Sprintf("token %s", token)
	}

	Version := "1"
	opts = append(opts,
		api.AddHeaderFunc("Authorization", getAuthValue),
		api.AddHeader("User-Agent", fmt.Sprintf("ghsettings %s", Version)),
		api.AddHeader("Accept", "application/vnd.github.antiope-preview+json"),
	)

	return api.NewClient(opts...), nil

}
