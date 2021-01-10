package prep

import (
	"log"

	"github.com/maseology/goHydro/grid"
	"github.com/maseology/goHydro/met"
	"github.com/maseology/mmio"
	"github.com/maseology/rdrr/model"
)

const metFP = "M:/ORMGP/met/ORMGP_5000_YCDB.met"

func loadMetFORC(toGD *grid.Definition) (*model.FORC, error) {

	h, c, _ := met.ReadMET(metFP, true)

	dat := make([][][]float64, 4) // [ 0:yield; 1:Ep ][staID][DateID]
	nloc, nstp := h.Nloc(), h.Nstep()
	for k := 0; k < 4; k++ {
		dat[k] = make([][]float64, nloc)
	}

	for i, a := range c.D { // dates
		for j, v := range a { // locations
			if i == 0 {
				for k := 0; k < 4; k++ {
					dat[k][j] = make([]float64, nstp)
				}
			}

			// tx, tn, r, s := v[0], v[1], v[2], v[3]
			for k := 0; k < 4; k++ {
				dat[k][j][i] = v[k]
			}
		}
	}

	mxr := make(map[int]int, toGD.Nact) // cell ID to met loc ID
	mgd := func() *grid.Definition {
		gd, err := grid.ReadGDEF(mmio.RemoveExtension(metFP)+".gdef", true)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if len(gd.Sactives) <= 0 {
			log.Fatalf("error: grid definition requires active cells")
		}
		return gd
	}()

	for _, cid := range toGD.Sactives {
		cntrd := toGD.Coord[cid]
		mxr[cid] = mgd.PointToCellID(cntrd.X, cntrd.Y)
	}

	return &model.FORC{
		T:  c.T,
		D:  dat,
		XR: mxr,
	}, nil
}
