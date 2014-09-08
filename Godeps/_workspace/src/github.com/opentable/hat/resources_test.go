package hat

import (
	"net/http"
	"testing"
)

func TestHat(t *testing.T) {
	if s, err := NewServer(Root{}); err != nil {
		t.Error(err)
	} else {
		println("Listening on :8080")
		t.Fatal(http.ListenAndServe(":8080", s))
	}
}

type Root struct {
	Hello  string
	Apps   *Apps   `hat:"embed(Name); page(1,1)"`
	Health *Health `hat:"embed(); link()"`
}

type Health struct {
	Hello     string `json:"hello"`
	TestField string `json:"testfield"`
}

type Apps map[string]App

// func (*Apps) Fields(_ *Root, _ string) ([]string, error) {
// 	return []string{"Name"}, nil
// }

func (entity *Apps) Page(number int, _ *Root, _ string) ([]string, error) {
	if number == 0 {
		number = 1
	}
	size := 0
	if size == 0 {
		size = 10
	}
	ids := []string{}
	for id, _ := range the_apps {
		ids = append(ids, id)
	}
	(*entity) = the_apps
	return ids, nil
}

// func (entity *Apps) PageFiltered(fields []string, number int, size int, _ *Root, _ string) ([]*App, error) {
// 	apps := []App{}
// 	for _, a := range the_apps {
// 		apps := append(apps, a)
// 	}
// 	return apps, nil
// }

// func (entity *Apps) Total(_ *Root, _ string) (int, error) {
// 	return 0, nil
// }

// func (entity *Apps) Manifest(_ *Root, _ string) error {
// 	(*entity) = Apps{
// 		ArbitraryField: "Hello",
// 		SomeOtherField: "Hi!",
// 	}
// 	return nil
// }

// func (*Apps) Page(_ *Root, _ string) ([]string, error) {
// 	ids := []string{}
// 	for k, _ := range the_apps {
// 		ids := append(ids, k)
// 	}
// 	return ids, nil
// }

type App struct {
	Name     string    `json:"name"`
	ID       string    `json:"id"`
	Versions *Versions `hat:"link()"`
}

func (entity *App) Manifest(parent *Apps, id string) error {
	if app, ok := (*parent)[id]; ok {
		(*entity) = app
	}
	return nil
}

type Versions map[string]Version

type Version struct {
	ID      string `json:"id"`
	Version string `json:"version`
	Date    string `json:"date"`
}

func (entity *Root) Manifest() error {
	(*entity) = Root{
		Hello: "Welcome to the test API.",
	}
	return nil
}

func (entity *Health) Manifest(_ *Root) error {
	(*entity) = Health{
		Hello: "Feelin' good!",
	}
	return nil
}

var the_apps = Apps{
	"test-app": App{
		"Test App", "test-app",
		&Versions{
			"0.0.1": Version{"test-app-v0-0-1", "0.0.1", "May 2013"},
			"0.0.2": Version{"test-app-v0-0-2", "0.0.2", "August 2014"},
		},
	},
	"other-app": App{
		"Other App", "other-app",
		&Versions{
			"0.1.0": Version{"other-app-v0-0-1", "0.1.0", "June 2014"},
			"0.4.0": Version{"other-app-v0-0-2", "0.4.0", "July 2014"},
		},
	},
}

func (entity *Versions) Page(page int, parent *App) ([]string, error) {
	ids := make([]string, len(*parent.Versions))
	i := 0
	for id, _ := range *parent.Versions {
		ids[i] = id
		i++
	}
	return ids, nil
}

func (entity *Version) Manifest(parent *Versions, id string) error {
	if version, ok := (*parent)[id]; ok {
		(*entity) = version
	}
	return nil
}
