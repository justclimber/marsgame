package changelog

import (
	"math"
	"strconv"
)

var total int
var optimized int

// Optimize intermediate changelog for individual object if they can be interpolated (have const diffs at least for 3 changes)
func (ch *ChangeLog) Optimize() {
	changelogLen := len(ch.ChangesByTimeLog)
	total += changelogLen
	if changelogLen < 3 {
		return
	}
	i1, i2, ok := lookupForSameDiff(changelogLen-1, changelogLen-2, ch)
	if ok {
		ch.cutInterpolableChanges(i1, i2)
		optimized += changelogLen - len(ch.ChangesByTimeLog)
	}
	//log.Printf("Optimization %f%%: %d/%d\n", float64(optimized/total), optimized, total)
}

func (ch *ChangeLog) cutInterpolableChanges(i1, i2 int) {
	ch.ChangesByTimeLog = append(ch.ChangesByTimeLog[:i2+1], ch.ChangesByTimeLog[i1:]...)
}

func lookupForSameDiff(tailIndex int, index int, ch *ChangeLog) (int, int, bool) {
	dl, dr, dcr := getDiff(tailIndex, index, ch)
	dl1, dr1, dcr1 := getDiff(tailIndex-1, index-1, ch)

	if checkAreDiffEqualZero(dl, dl1) && checkAreDiffTheSame(dr, dr1) && checkAreDiffTheSame(dcr, dcr1) {
		if index > 2 {
			_, index1, ok := lookupForSameDiff(tailIndex, index-2, ch)
			if ok {
				return tailIndex, index1, ok
			}
			return 0, 0, false
		}
		return tailIndex, index - 1, true
	}
	return 0, 0, false
}

const floatDelta = 0.01

func checkAreDiffTheSame(d1, d2 map[string]*float64) bool {
	for k, v := range d1 {
		if d2[k] == nil || v == nil || *d2[k]-*v > floatDelta {
			return false
		}
	}
	return true
}

func checkAreDiffEqualZero(d1, d2 map[string]*float64) bool {
	for k, v := range d1 {
		if d2[k] == nil || v == nil || math.Abs(*d2[k]) > floatDelta || math.Abs(*v) > floatDelta {
			return false
		}
	}
	return true
}

func getDiff(index int, index1 int, ch *ChangeLog) (map[string]*float64, map[string]*float64, map[string]*float64) {
	l, r, cr := getValuesForChanges(ch.ChangesByTimeLog[index])
	l1, r1, cr1 := getValuesForChanges(ch.ChangesByTimeLog[index1])
	dl := getDiffFor2Maps(l, l1)
	dr := getDiffFor2Maps(r, r1)
	dcr := getDiffFor2Maps(cr, cr1)

	return dl, dr, dcr
}

func getDiffFor2Maps(d1, d2 map[string]*float64) map[string]*float64 {
	result := make(map[string]*float64)
	for k, v := range d1 {
		if v != nil && d2[k] != nil {
			result[k] = new(float64)
			d := *v - *d2[k]
			result[k] = &d
		} else {
			result[k] = nil
		}
	}
	return result
}

func getValuesForChanges(changeByTime *ChangeByTime) (map[string]*float64, map[string]*float64, map[string]*float64) {
	l := make(map[string]*float64)
	r := make(map[string]*float64)
	cr := make(map[string]*float64)
	for _, changeByObject := range changeByTime.ChangesByObject {
		key := strconv.Itoa(changeByObject.ObjId) + " " + changeByObject.ObjType
		l[key] = changeByObject.Length
		r[key] = changeByObject.Angle
		cr[key] = changeByObject.CannonAngle
	}
	return l, r, cr
}
