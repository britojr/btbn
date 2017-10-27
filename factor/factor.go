package factor

import (
	"math/rand"
	"time"

	"github.com/britojr/btbn/vars"
	"gonum.org/v1/gonum/floats"
)

// Factor defines a function that maps a joint of categorical variables to float values
// every factor operation, if not explicitly stated otherwise, modifies the given factor
type Factor struct {
	values []float64
	vs     vars.VarList
}

var (
	opAdd = func(a, b float64) float64 { return a + b }
	opSub = func(a, b float64) float64 { return a - b }
	opDiv = func(a, b float64) float64 {
		if b == 0 {
			panic("factor: op division by zero")
		}
		return a / b
	}
	opMul = func(a, b float64) float64 { return a * b }
)

var seed = func() int64 {
	return time.Now().UnixNano()
}

// New creates a new factor with uniform distribution for given variables
// a factor with no variables have a value of one
func New(xs ...*vars.Var) (f *Factor) {
	f = new(Factor)
	f.vs = vars.VarList(xs).Copy()
	f.values = make([]float64, f.vs.NStates())
	tot := float64(len(f.values))
	for i := range f.values {
		f.values[i] = 1 / tot
	}
	return
}

// Copy returns a copy of f
func (f *Factor) Copy() (g *Factor) {
	g = new(Factor)
	g.vs = f.vs.Copy()
	g.values = make([]float64, len(f.values))
	copy(g.values, f.values)
	return
}

// SetValues sets the given values to the factor
func (f *Factor) SetValues(values []float64) *Factor {
	copy(f.values, values)
	return f
}

// Values return reference for factor values
func (f *Factor) Values() []float64 {
	return f.values
}

// Variables return reference for factor variables
func (f *Factor) Variables() vars.VarList {
	return f.vs
}

// RandomDistribute sets values with a random distribution
func (f *Factor) RandomDistribute(xs ...*vars.Var) *Factor {
	r := rand.New(rand.NewSource(seed()))
	for i := range f.values {
		f.values[i] = r.Float64()
		for f.values[i] <= 0 {
			f.values[i] = r.Float64()
		}
	}
	return f.Normalize(xs...)
}

// Plus adds g to f
func (f *Factor) Plus(g *Factor) *Factor {
	return f.operationIn(g, opAdd)
}

// Times multiplies f by g
func (f *Factor) Times(g *Factor) *Factor {
	return f.operationIn(g, opMul)
}

// operation applies given operation as f = f op g
func (f *Factor) operationIn(g *Factor, op func(a, b float64) float64) *Factor {
	if f.vs.Equal(g.vs) {
		return f.operationEq(g, op)
	}
	return f.operationTr(f.Copy(), g, op)
}

// operationEq applies given operation as f = f op g assuming equal sets of variables
func (f *Factor) operationEq(g *Factor, op func(a, b float64) float64) *Factor {
	for i, v := range g.values {
		f.values[i] = op(f.values[i], v)
	}
	return f
}

// operationTr returns applies given operation as operation f = f1 op f2
func (f *Factor) operationTr(f1, f2 *Factor, op func(a, b float64) float64) *Factor {
	f.vs = f1.vs.Union(f2.vs)
	f.values = make([]float64, f.vs.NStates())
	ixf1 := vars.NewIndexFor(f1.vs, f.vs)
	ixf2 := vars.NewIndexFor(f2.vs, f.vs)
	for i := range f.values {
		f.values[i] = op(f1.values[ixf1.I()], f2.values[ixf2.I()])
		ixf1.Next()
		ixf2.Next()
	}
	return f
}

// Normalize normalizes f so the values add up to one
// the normalization can be conditional in one variable
func (f *Factor) Normalize(xs ...*vars.Var) *Factor {
	if len(xs) == 0 {
		return f.normalizeAll()
	}
	condVars := f.vs.Diff(xs)
	// if len(condVars) == 0 {
	// 	return f.normalizeAll()
	// }
	ixf := vars.NewIndexFor(condVars, f.vs)
	sums := make([]float64, condVars.NStates())
	for _, v := range f.values {
		sums[ixf.I()] += v
		ixf.Next()
	}
	ixf.Reset()
	for i := range f.values {
		if sums[ixf.I()] != 0 {
			f.values[i] /= sums[ixf.I()]
		} else {
			panic("factor: conditional prob with zero sum")
		}
		ixf.Next()
	}
	return f
}

func (f *Factor) normalizeAll() *Factor {
	sum := floats.Sum(f.values)
	if sum != 0 {
		for i := range f.values {
			f.values[i] /= sum
		}
	} else {
		panic("factor: all values add up to zero")
	}
	return f
}

// SumOut sums out the given variables
func (f *Factor) SumOut(xs ...*vars.Var) *Factor {
	if len(xs) == 0 {
		return f
	}
	oldVs := f.vs
	oldVal := f.values
	f.vs = oldVs.Diff(xs)
	f.values = make([]float64, f.vs.NStates())
	ixf := vars.NewIndexFor(f.vs, oldVs)
	for _, v := range oldVal {
		f.values[ixf.I()] += v
		ixf.Next()
	}
	return f
}

// Marginalize projects the distribution in the given variables
func (f *Factor) Marginalize(xs ...*vars.Var) *Factor {
	return f.SumOut(f.vs.Diff(xs)...)
}

// Reduce silences the values that are not compatible with the given evidence
func (f *Factor) Reduce(e map[int]int) *Factor {
	panic("factor: not implemented")
}

// TimesNew creates a new factor h = f * g
// func (f *Factor) TimesNew(g *Factor) *Factor {
// 	return f.operationNew(g, opMul)
// }
// operationNew applies given operation resulting in new factor h = f op g
// func (f *Factor) operationNew(g *Factor, op func(a, b float64) float64) *Factor {
// 	if f.vs.Equal(g.vs) {
// 		return f.Copy().operationEq(g, op)
// 	}
// 	return f.Copy().operationTr(f, g, op)
// }
// // MarginalizeNew returns a new factor summing out the variables not given
// func (f *Factor) MarginalizeNew(xs ...*vars.Var) *Factor {
// 	panic("factor: not implemented")
// }
// SumOutNew returns a new factor with the given variables summed out
// func (f *Factor) SumOutNew(xs ...*vars.Var) *Factor {
// 	return f.Copy().SumOut(xs...)
// }
// SumOutID sums out the variables given by id
// func (f *Factor) SumOutID(ids ...int) *Factor {
// 	return f.SumOut(f.vs.IntersecID(ids...)...)
// }
//
// // SumOutIDNew returns a new factor with the given variables summed out
// func (f *Factor) SumOutIDNew(ids ...int) *Factor {
// 	return f.SumOutNew(f.vs.IntersecID(ids...)...)
// }
