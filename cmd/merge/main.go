package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/test-infra/prow/config/org"
)

func parseKeyValue(s string) (string, string) {
	p := strings.SplitN(s, "=", 2)
	if len(p) == 1 {
		return p[0], ""
	}
	return p[0], p[1]
}

type orgFlagMap map[string]string

func (fm orgFlagMap) String() string {
	var parts []string
	for key, value := range fm {
		if value == "" {
			continue
		}
		parts = append(parts, key+"="+value)
	}
	return strings.Join(parts, ",")
}

func (fm orgFlagMap) Set(s string) error {
	k, v := parseKeyValue(s)
	if _, present := fm[k]; present {
		return fmt.Errorf("duplicate key: %s", k)
	}
	fm[k] = v
	return nil
}

type options struct {
	orgs        orgFlagMap
	mergeTeams  bool
	ignoreTeams bool
}

func main() {
	o := options{orgs: orgFlagMap{}}
	flag.Var(o.orgs, "org-part", "Each instance adds an org-name=org.yaml part")
	flag.BoolVar(&o.mergeTeams, "merge-teams", false, "Merge team-name/team.yaml files in each org.yaml dir")
	flag.BoolVar(&o.ignoreTeams, "ignore-teams", false, "Never configure teams")
	flag.Parse()

	for _, a := range flag.Args() {
		logrus.Print("Extra", a)
		o.orgs.Set(a)
	}

	if o.mergeTeams && o.ignoreTeams {
		logrus.Fatal("--merge-teams xor --ignore-teams, not both")
	}

	cfg, err := loadOrgs(o)
	if err != nil {
		logrus.Fatalf("Failed to load orgs: %v", err)
	}
	pc := org.FullConfig{
		Orgs: cfg,
	}
	out, err := yaml.Marshal(pc)
	if err != nil {
		logrus.Fatalf("Failed to marshal orgs: %v", err)
	}
	fmt.Println(string(out))
}

func unmarshal(path string) (*org.Config, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read: %v", err)
	}
	var cfg org.Config
	if err := yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	return &cfg, nil
}

func loadOrgs(o options) (map[string]org.Config, error) {
	config := map[string]org.Config{}
	for name, path := range o.orgs {
		cfg, err := unmarshal(path)
		if err != nil {
			return nil, fmt.Errorf("error in %s: %v", path, err)
		}
		switch {
		case o.ignoreTeams:
			cfg.Teams = nil
		case o.mergeTeams:
			if cfg.Teams == nil {
				cfg.Teams = map[string]org.Team{}
			}
			prefix := filepath.Dir(path)
			err := filepath.Walk(prefix, func(path string, info os.FileInfo, err error) error {
				switch {
				case path == prefix:
					return nil // Skip base dir
				case info.IsDir() && filepath.Dir(path) != prefix:
					logrus.Infof("Skipping %s and its children", path)
					return filepath.SkipDir // Skip prefix/foo/bar/ dirs
				case !info.IsDir() && filepath.Dir(path) == prefix:
					return nil // Ignore prefix/foo files
				case filepath.Base(path) == "teams.yaml":
					teamCfg, err := unmarshal(path)
					if err != nil {
						return fmt.Errorf("error in %s: %v", path, err)
					}

					for name, team := range teamCfg.Teams {
						cfg.Teams[name] = team
					}
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("merge teams %s: %v", path, err)
			}
		}
		config[name] = *cfg
	}
	return config, nil
}
