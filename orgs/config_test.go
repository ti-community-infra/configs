package orgs

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/test-infra/prow/config/org"
	"k8s.io/test-infra/prow/github"

	"github.com/ghodss/yaml"
)

const botName = "ti-community-prow-bot"

func testDuplicates(list sets.String) error {
	found := sets.String{}
	dups := sets.String{}
	all := list.List()
	for _, i := range all {
		if found.Has(i) {
			dups.Insert(i)
		}
		found.Insert(i)
	}
	if n := len(dups); n > 0 {
		return fmt.Errorf("%d duplicate names: %s", n, strings.Join(dups.List(), ", "))
	}
	return nil
}

func sortList(list []string) []string {
	items := make([]string, len(list))
	for _, l := range list {
		items = append(items, strings.ToLower(l))
	}
	sort.Strings(items)
	return items
}

func isSorted(list []string) bool {
	items := make([]string, len(list))
	for _, l := range list {
		items = append(items, strings.ToLower(l))
	}

	return sort.StringsAreSorted(items)
}

func normalize(s sets.String) sets.String {
	out := sets.String{}
	for i := range s {
		out.Insert(github.NormLogin(i))
	}
	return out
}

// testTeamMembers ensures that a user is not a maintainer and member at the same time,
// there are no duplicate names in the list and all users are org members.
func testTeamMembers(teams map[string]org.Team, admins sets.String, orgMembers sets.String, orgName string) []error {
	var errs []error
	for teamName, team := range teams {
		teamMaintainers := sets.NewString(team.Maintainers...)
		teamMembers := sets.NewString(team.Members...)

		teamMaintainers = normalize(teamMaintainers)
		teamMembers = normalize(teamMembers)

		// ensure all teams have privacy as closed
		if team.Privacy == nil || (team.Privacy != nil && *team.Privacy != org.Closed) {
			errs = append(errs, fmt.Errorf("The team %s in org %s doesn't have the `privacy: closed` field", teamName, orgName))
		}

		// check for non-admins in maintainers list
		if nonAdminMaintainers := teamMaintainers.Difference(admins); len(nonAdminMaintainers) > 0 {
			errs = append(errs, fmt.Errorf("The team %s in org %s has non-admins listed as maintainers; these users should be in the members list instead: %s", teamName, orgName, strings.Join(nonAdminMaintainers.List(), ",")))
		}

		// check for users in both maintainers and members
		if both := teamMaintainers.Intersection(teamMembers); len(both) > 0 {
			errs = append(errs, fmt.Errorf("The team %s in org %s has users in both maintainer admin and member roles: %s", teamName, orgName, strings.Join(both.List(), ", ")))
		}

		// check for duplicates
		if err := testDuplicates(teamMaintainers); err != nil {
			errs = append(errs, fmt.Errorf("The team %s in org %s has duplicate maintainers: %v", teamName, orgName, err))
		}
		if err := testDuplicates(teamMembers); err != nil {
			errs = append(errs, fmt.Errorf("The team %s in org %s has duplicate members: %v", teamMembers, orgName, err))
		}

		// check if all are org members
		if missing := teamMembers.Difference(orgMembers); len(missing) > 0 {
			errs = append(errs, fmt.Errorf("The following members of team %s are not %s org members: %s", teamName, orgName, strings.Join(missing.List(), ", ")))
		}

		// check if admins are a regular member of team
		if adminTeamMembers := teamMembers.Intersection(admins); len(adminTeamMembers) > 0 {
			errs = append(errs, fmt.Errorf("The team %s in org %s has org admins listed as members; these users should be in the maintainers list instead, and cannot be on the members list: %s", teamName, orgName, strings.Join(adminTeamMembers.List(), ", ")))
		}

		// check if lists are sorted
		if !isSorted(team.Maintainers) {
			errs = append(errs, fmt.Errorf("the team %s in org %s has an unsorted list of maintainers; it should be %v", teamName, orgName, sortList(team.Maintainers)))
		}
		if !isSorted(team.Members) {
			errs = append(errs, fmt.Errorf("the team %s in org %s has an unsorted list of members; it should be %v", teamName, orgName, sortList(team.Members)))
		}

		if team.Children != nil {
			errs = append(errs, testTeamMembers(team.Children, admins, orgMembers, orgName)...)
		}
	}
	return errs
}

func TestAllOrgs(t *testing.T) {
	raw, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		t.Fatalf("cannot read config: %v", err)
	}
	var cfg org.FullConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		t.Fatalf("cannot read config.yaml from //config:gen-config.yaml: %v", err)
	}

	for _, org := range cfg.Orgs {
		members := normalize(sets.NewString(org.Members...))
		admins := normalize(sets.NewString(org.Admins...))
		allOrgMembers := members.Union(admins)

		if both := admins.Intersection(members); len(both) > 0 {
			t.Errorf("users in both org admin and member roles for org '%s': %s", *org.Name, strings.Join(both.List(), ", "))
		}

		if !admins.Has(botName) {
			t.Errorf(botName + " must be an admin")
		}

		if err := testDuplicates(admins); err != nil {
			t.Errorf("duplicate admins: %v", err)
		}
		if err := testDuplicates(allOrgMembers); err != nil {
			t.Errorf("duplicate members: %v", err)
		}
		if !isSorted(org.Admins) {
			t.Errorf("admins for %s org are unsorted; it should be %v", *org.Name, sortList(org.Admins))
		}
		if !isSorted(org.Members) {
			t.Errorf("members for %s org are unsorted; it should be %v", *org.Name, sortList(org.Members))
		}

		if errs := testTeamMembers(org.Teams, admins, allOrgMembers, *org.Name); errs != nil {
			for _, err := range errs {
				t.Error(err)
			}
		}
	}
}
