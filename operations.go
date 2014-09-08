package main

import (
	"github.com/opentable/hat"
)

var root = &Root{Hello: "Deployment API; Submit bugs & feature requests to https://github.com/opentable/deploy/issues"}

func (r *Root) Manifest() error {
	*r = *root
	return nil
}

func (p *Pools) Page() ([]string, error) {
	return state.GetPoolIDs(), nil
}

func (p *Pool) Manifest(pools *Pools, name string) error {
	println("Manifest pool "+name+" got ", p)
	pool := state.GetPool(name)
	if pool != nil {
		*p = *pool
	}
	return nil
}

func (p *Pool) Write(_ *Pools, name string) error {
	pool := state.GetPool(name)
	if pool != nil {
		return hat.HttpError(409, "Pool "+name+" already exists.")
	}
	state.SetPool(name, p)
	return nil
}

func (a *Apps) Page() ([]string, error) {
	return []string{}, nil
}

func (a *App) Manifest(_ *Apps, name string) error {
	return nil
}

func (a *App) Write(as *Apps, name string) error {
	(*as)[name] = a
	return nil
}

func (v *Versions) Page() ([]string, error) {
	return []string{}, nil
}

func (v *Version) Manifest(_ *Versions, name string) error {
	return nil
}

func (v *Version) Write(vs *Versions, name string) error {
	if _, exists := (*vs)[name]; exists {
		return hat.HttpError(409, name+" already exists.")
	}
	*(*vs)[name] = *v
	return nil
}
