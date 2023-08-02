package helper

import "fmt"

type HRUnitLevel int

const (
	HRUnitLevelAuto HRUnitLevel = iota
	HRUnitLevelB
	HRUnitLevelKB
	HRUnitLevelMB
	HRUnitLevelGB
	HRUnitLevelTB
	HRUnitLevelPB
)

func CalcVUnitLevel(v float64, base uint) HRUnitLevel {
	var level HRUnitLevel

	fBase := float64(base)

	if v < fBase {
		level = HRUnitLevelB
	} else if v < fBase*fBase {
		level = HRUnitLevelKB
	} else if v < fBase*fBase*fBase {
		level = HRUnitLevelMB
	} else if v < fBase*fBase*fBase*fBase {
		level = HRUnitLevelGB
	} else if v < fBase*fBase*fBase*fBase*fBase {
		level = HRUnitLevelTB
	} else {
		level = HRUnitLevelPB
	}

	return level
}

func FormatVEx(v float64, base uint, u string, level HRUnitLevel) (float64, string) {
	if level == HRUnitLevelAuto {
		level = CalcVUnitLevel(v, base)
	}

	fBase := float64(base)

	switch level {
	default:
		fallthrough
	case HRUnitLevelB:
		return v, fmt.Sprintf("%s", u)
	case HRUnitLevelKB:
		return v / fBase, fmt.Sprintf("K%s", u)
	case HRUnitLevelMB:
		return v / fBase / fBase, fmt.Sprintf("M%s", u)
	case HRUnitLevelGB:
		return v / fBase / fBase / fBase, fmt.Sprintf("G%s", u)
	case HRUnitLevelTB:
		return v / fBase / fBase / fBase / fBase, fmt.Sprintf("T%s", u)
	case HRUnitLevelPB:
		return v / fBase / fBase / fBase / fBase / fBase, fmt.Sprintf("P%s", u)
	}
}

func FormatVExToString(v float64, base uint, u string, level HRUnitLevel) string {
	if level == HRUnitLevelAuto {
		level = CalcVUnitLevel(v, base)
	}

	v, s := FormatVEx(v, base, u, level)
	return fmt.Sprintf("%.02f %s", v, s)
}
