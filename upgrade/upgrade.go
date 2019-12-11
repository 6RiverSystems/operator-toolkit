package upgrade

// Interface is an interface to be implemented for ops that perform upgrade
// of resources by provided specifications.
type Interface interface {
	ResourceLen() int
	SpecLen() int
	SpecBelongsToResource(spec, resource int) bool
	ResourceNeedsUpdate(resource, spec int) bool
	CreateResource(spec int) error
	UpdateResource(resource, spec int) error
	DeleteResource(resource int) error
}

// Upgrade upgrades resources according to provided specifications.
// It also prunes resources that are not longer founnd among specification list.
func Upgrade(data Interface) error {
	specToRes := make(map[int]int)
	resourceToSpec := make(map[int]int)

	for spec := 0; spec < data.SpecLen(); spec++ {
		for resource := 0; resource < data.ResourceLen(); resource++ {
			if data.SpecBelongsToResource(spec, resource) {
				specToRes[spec] = resource
				resourceToSpec[resource] = spec
				break
			}
		}
	}

	for spec := 0; spec < data.SpecLen(); spec++ {
		if resource, ok := specToRes[spec]; ok {
			if data.ResourceNeedsUpdate(resource, spec) {
				if err := data.UpdateResource(resource, spec); err != nil {
					return err
				}
			}
		} else {
			if err := data.CreateResource(spec); err != nil {
				return err
			}
		}
	}

	for resource := 0; resource < data.ResourceLen(); resource++ {
		if _, ok := resourceToSpec[resource]; !ok {
			if err := data.DeleteResource(resource); err != nil {
				return err
			}
		}
	}

	return nil
}
