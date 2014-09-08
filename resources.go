package main

type Tags struct {
	Tags []string `json:"tags"`
}

type Root struct {
	Hello string `json:"hello"`
	Pools *Pools `hat:"embed()"`
}

type Pools map[string]*Pool

type Pool struct {
	Name         string            `json:"name"`
	MarathonHost string            `json:"marathonHost"`
	Env          map[string]string `json:"env"`
	Apps         *Apps             `hat:"embed()"`
	Tags
}

type Apps map[string]*App

type App struct {
	Name     string `json:"name"`
	Versions *Versions
	Tags
}

type Versions map[string]*Version

type Version struct {
	AppName      string       `json:"appName"`
	Version      string       `json:"version"`
	ArtifactURLs []string     `json:"artifactUrls"`
	Command      []string     `json:"command"`
	HealthURI    string       `json:"healthUri"`
	MinInstances int          `json:"minInstances"`
	MaxInstances int          `json:"maxInstances"`
	Requirements Requirements `json:"requirements"`
	Tags
}

type Requirements struct {
	Ports         int
	SpecificPorts []int
	CPU           float64
	MemoryMB      float64
	DiskMB        float64
}
