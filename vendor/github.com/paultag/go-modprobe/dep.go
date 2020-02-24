package modprobe

import (
	"bufio"
	"os"
	"strings"

	"pault.ag/go/topsort"
)

// Given a path to a .ko file, determine what modules will have to be present
// before loading that module.
func Dependencies(path string) ([]string, error) {
	deps, err := loadDependencies()
	if err != nil {
		return nil, err
	}
	return deps.Load(path)
}

// simple container type that stores a mapping from an element to elements
// that it depends on.
type dependencies map[string][]string

// top level loading of the dependency tree. this will start a network
// walk the dep tree, load them into the network, and return a topological
// sort of the modules.
func (d dependencies) Load(name string) ([]string, error) {
	network := topsort.NewNetwork()
	if err := d.load(name, network); err != nil {
		return nil, err
	}

	order, err := network.Sort()
	if err != nil {
		return nil, err
	}

	ret := []string{}
	for _, node := range order {
		ret = append(ret, node.Name)
	}
	return ret, nil
}

// add a specific dependency to the network, and recurse on the leafs.
func (d dependencies) load(name string, network *topsort.Network) error {
	if network.Get(name) != nil {
		return nil
	}
	network.AddNode(name, nil)

	for _, dep := range d[name] {
		if err := d.load(dep, network); err != nil {
			return err
		}
		if err := network.AddEdge(dep, name); err != nil {
			return err
		}
	}

	return nil
}

// get a dependency map from the running kernel's modules.dep file
func loadDependencies() (dependencies, error) {
	path := modulePath("modules.dep")

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	deps := map[string][]string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		chunks := strings.SplitN(scanner.Text(), ":", 2)
		depString := strings.TrimSpace(chunks[1])
		if len(depString) == 0 {
			continue
		}

		ret := []string{}
		for _, dep := range strings.Split(depString, " ") {
			ret = append(ret, modulePath(dep))
		}
		deps[modulePath(chunks[0])] = ret
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return deps, nil
}
