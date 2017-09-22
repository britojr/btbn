package optimizer

import (
	"log"

	"github.com/britojr/btbn/scr"
	"github.com/britojr/utl/conv"
)

// common defines optimizer's commom/default behaviours
type common struct {
	tw           int          // treewidth
	nv           int          // number of variables
	scoreRankers []scr.Ranker // score rankers for each variable
}

func newCommon(scoreRankers []scr.Ranker) *common {
	s := new(common)
	s.scoreRankers = scoreRankers
	s.nv = len(s.scoreRankers)
	return s
}

func (s *common) SetDefaultParameters() {
	s.tw = 3
}

func (s *common) SetFileParameters(parms map[string]string) {
	if tw, ok := parms[ParmTreewidth]; ok {
		s.tw = conv.Atoi(tw)
	}
}

func (s *common) ValidateParameters() {
	if s.tw <= 0 || s.nv < s.tw+2 {
		log.Printf("n=%v, tw=%v\n", s.nv, s.tw)
		log.Panic("Invalid treewidth! Choose values such that: n >= tw+2 and tw > 0")
	}
}

func (s *common) PrintParameters() {
	log.Printf(" ========== ALGORITHM PARAMETERS ========== \n")
	log.Printf("number of variables: %v\n", s.nv)
	log.Printf("%v: %v\n", ParmTreewidth, s.tw)
}
