package main

import (
	log "github.com/Sirupsen/logrus"

	"github.com/mkrakowitzer/ghsettings/command"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func main() {
	command.Execute()
}
