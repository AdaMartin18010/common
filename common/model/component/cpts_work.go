package component

import (
	multierror "github.com/hashicorp/go-multierror"
)

var (
	_ CptsOperator = (*Cpts)(nil)
	//_ Component          = (*Components)(nil)
)

// Cpts is a collection of Component
type Cpts []Cpt

// v is nil value append also  means that len > 0
func NewCpts(v ...Cpt) Cpts {
	return append([]Cpt{}, v...)
}

// Len returns the amount of Components in the collection.
func (cps *Cpts) Len() int {
	return len(*cps)
}

func (cps *Cpts) AddCpts(cmpts ...Cpt) {
	*cps = append(([]Cpt)(*cps), cmpts...)
}

func (cps *Cpts) RemoveCpts(cmpts ...Cpt) {
	if len(cmpts) == 0 {
		return
	}

	for i := 0; i < len(*cps); i++ {
		cmp := ([]Cpt)(*cps)[i]
		if cmp == nil {
			*cps = append(([]Cpt)(*cps)[:i], ([]Cpt)(*cps)[i+1:]...)
			i--
			continue
		}

		for j := 0; j < len(cmpts); j++ {
			if cmp == cmpts[j] {
				*cps = append(([]Cpt)(*cps)[:i], ([]Cpt)(*cps)[i+1:]...)
				i--
			}
		}
	}
}

func (cps *Cpts) Cpt(idname IdName) Cpt {
	for _, cp := range *cps {
		if cp != nil && cp.Id() == idname {
			return cp
		}
	}

	return nil
}

// Each enumerates through the Components and calls specified callback function.
func (cps *Cpts) Each(f func(Cpt)) {
	for _, cp := range *cps {
		if cp != nil {
			f(cp)
		}
	}
}

// Start calls the Start method of each Component in the collection
func (cps *Cpts) Start() (err error) {
	for _, cp := range *cps {
		if cp == nil {
			continue
		}

		if cp.IsRunning() {
			continue
		}

		if rerr := cp.Start(); rerr != nil {
			err = multierror.Append(err, rerr)
		}
	}

	return
}

// Stop calls the Stop method of each Component in the collection
func (cps *Cpts) Stop() (err error) {
	for _, cp := range *cps {
		if cp == nil {
			continue
		}

		if !cp.IsRunning() {
			continue
		}

		if rerr := cp.Stop(); rerr != nil {
			err = multierror.Append(err, rerr)
		}
	}

	return
}
