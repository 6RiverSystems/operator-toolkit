package upgrade_test

import (
	"testing"

	. "github.com/6RiverSystems/operator-toolkit/upgrade"
	"github.com/go-test/deep"
)

type spec struct {
	name string
}

type upgradeTester struct {
	resources   []string
	specs       []spec
	createCalls []int
	updateCalls [][2]int
	deleteCalls []int
}

func newUpgradeTester(resources []string, specs []spec) *upgradeTester {
	return &upgradeTester{
		resources:   resources,
		specs:       specs,
		createCalls: make([]int, 0),
		updateCalls: make([][2]int, 0),
		deleteCalls: make([]int, 0),
	}
}

func (u *upgradeTester) ResourceNeedsUpdate(resource, spec int) bool {
	return true
}

func (u *upgradeTester) ResourceLen() int { return len(u.resources) }

func (u *upgradeTester) SpecLen() int { return len(u.specs) }

func (u *upgradeTester) SpecBelongsToResource(spec, resource int) bool {
	return u.specs[spec].name == u.resources[resource]
}

func (u *upgradeTester) CreateResource(spec int) error {
	u.createCalls = append(u.createCalls, spec)
	return nil
}

func (u *upgradeTester) UpdateResource(resource, spec int) error {
	u.updateCalls = append(u.updateCalls, [2]int{resource, spec})
	return nil
}

func (u *upgradeTester) DeleteResource(resource int) error {
	u.deleteCalls = append(u.deleteCalls, resource)
	return nil
}

func TestUpgrade(t *testing.T) {
	tester := newUpgradeTester(
		[]string{
			"vodka",
			"gin",
			"beer",
		},
		[]spec{
			spec{name: "gin"},
			spec{name: "beer"},
			spec{name: "ron"},
		},
	)

	if err := Upgrade(tester); err != nil {
		t.Error(err)
	}
	if diff := deep.Equal([]int{2}, tester.createCalls); diff != nil {
		t.Error(diff)
	}
	if diff := deep.Equal([][2]int{[2]int{1, 0}, [2]int{2, 1}}, tester.updateCalls); diff != nil {
		t.Error(diff)
	}
	if diff := deep.Equal([]int{0}, tester.deleteCalls); diff != nil {
		t.Error(diff)
	}
}

func TestDeleteAll(t *testing.T) {
	tester := newUpgradeTester(
		[]string{
			"vodka",
			"gin",
			"beer",
		},
		nil,
	)

	if err := Upgrade(tester); err != nil {
		t.Error(err)
	}
	if diff := deep.Equal([]int{}, tester.createCalls); diff != nil {
		t.Error(diff)
	}
	if diff := deep.Equal([][2]int{}, tester.updateCalls); diff != nil {
		t.Error(diff)
	}
	if diff := deep.Equal([]int{0, 1, 2}, tester.deleteCalls); diff != nil {
		t.Error(diff)
	}
}

func TestCreateAll(t *testing.T) {
	tester := newUpgradeTester(
		nil,
		[]spec{
			spec{name: "gin"},
			spec{name: "beer"},
			spec{name: "ron"},
		},
	)

	if err := Upgrade(tester); err != nil {
		t.Error(err)
	}
	if diff := deep.Equal([]int{0, 1, 2}, tester.createCalls); diff != nil {
		t.Error(diff)
	}
	if diff := deep.Equal([][2]int{}, tester.updateCalls); diff != nil {
		t.Error(diff)
	}
	if diff := deep.Equal([]int{}, tester.deleteCalls); diff != nil {
		t.Error(diff)
	}
}
