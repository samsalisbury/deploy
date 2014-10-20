package main

import (
	"github.com/opentable/hat"
	"github.com/opentable/ot-go-lib/env"
	"github.com/opentable/ot-go-lib/gutter"
	"github.com/opentable/ot-go-lib/logging"
	"github.com/opentable/ot-go-lib/service"
)

var (
	stateGitURL  = env.RequireString("OT_DEPLOY_STATE_REPO_URL")
	configGitURL = env.RequireString("OT_CLOUD_PLATFORM_CONFIG_REPO")
	log          = logging.StandardConfig("deploy").StartupLog(0)
	gitState     = initGitClient(stateGitURL)
	gitConfig    = initGitClient(configGitURL)
	svc          service.Service
)

func main() {
	s, err := hat.NewServer(Root{})
	if err != nil {
		log.Fatal(err)
		return
	}
	svc, err := service.NewHTTPServiceFromEnv("deploy", s.ServeHTTP)
	if err != nil {
		log.Fatal(err)
		return
	}
	svc.Start()
}

func initGitClient(url string) gutter.Client {
	git, err := gutter.NewGitClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := git.Clone(url); err != nil {
		log.Fatal(err)
	}
	return git
}
