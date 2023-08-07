package helper

import "fmt"

type VUnitLevel int

const (
	VUnitLevelAuto VUnitLevel = iota
	VUnitLevelB
	VUnitLevelKB
	VUnitLevelMB
	VUnitLevelGB
	VUnitLevelTB
	VUnitLevelPB
)

func CalcVUnitLevel(v float64, base uint) VUnitLevel {
	var level VUnitLevel

	fBase := float64(base)

	if v < fBase {
		level = VUnitLevelB
	} else if v < fBase*fBase {
		level = VUnitLevelKB
	} else if v < fBase*fBase*fBase {
		level = VUnitLevelMB
	} else if v < fBase*fBase*fBase*fBase {
		level = VUnitLevelGB
	} else if v < fBase*fBase*fBase*fBase*fBase {
		level = VUnitLevelTB
	} else {
		level = VUnitLevelPB
	}

	return level
}

func FormatVEx(v float64, base uint, u string, level VUnitLevel) (float64, string) {
	if level == VUnitLevelAuto {
		level = CalcVUnitLevel(v, base)
	}

	fBase := float64(base)

	switch level {
	default:
		fallthrough
	case VUnitLevelB:
		return v, u
	case VUnitLevelKB:
		return v / fBase, fmt.Sprintf("K%s", u)
	case VUnitLevelMB:
		return v / fBase / fBase, fmt.Sprintf("M%s", u)
	case VUnitLevelGB:
		return v / fBase / fBase / fBase, fmt.Sprintf("G%s", u)
	case VUnitLevelTB:
		return v / fBase / fBase / fBase / fBase, fmt.Sprintf("T%s", u)
	case VUnitLevelPB:
		return v / fBase / fBase / fBase / fBase / fBase, fmt.Sprintf("P%s", u)
	}
}

func FormatVExToString(v float64, base uint, u string, level VUnitLevel) string {
	if level == VUnitLevelAuto {
		level = CalcVUnitLevel(v, base)
	}

	v, s := FormatVEx(v, base, u, level)

	return fmt.Sprintf("%.02f %s", v, s)
}
