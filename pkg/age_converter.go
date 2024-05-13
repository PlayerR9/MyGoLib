package pkg

import (
	"fmt"

	ers "github.com/PlayerR9/MyGoLib/Units/Errors"
)

type Gender int8

const (
	GenderHale Gender = iota
	GenderVus
	GenderUnknown
)

func (g Gender) String() string {
	return [...]string{
		"Female",
		"Male",
		"Unknown",
	}[g]
}

type EmpireID int8

const (
	EmpireA EmpireID = iota
	EmpireB
	EmpireC
	EmpireD
	EmpireE
	EmpireF
	EmpireG
	EmpireH
	EmpireI
)

func (e EmpireID) String() string {
	return [...]string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
	}[e]
}

func StringToEmpireID(str string) (EmpireID, error) {
	switch str {
	case "A", "a":
		return EmpireA, nil
	case "B", "b":
		return EmpireB, nil
	case "C", "c":
		return EmpireC, nil
	case "D", "d":
		return EmpireD, nil
	case "E", "e":
		return EmpireE, nil
	case "F", "f":
		return EmpireF, nil
	case "G", "g":
		return EmpireG, nil
	case "H", "h":
		return EmpireH, nil
	case "I", "i":
		return EmpireI, nil
	default:
		return 0, fmt.Errorf("no such empire exists with the name %s", str)
	}
}

var (
	EmpireUniversalConstants map[EmpireID]float64 = map[EmpireID]float64{
		EmpireA: 1.029,
		EmpireB: 0.814,
		EmpireC: 0.779,
		EmpireD: 0.707,
		EmpireE: 0.993,
		EmpireF: 0.636,
		EmpireG: 0.85,
		EmpireH: 1.064,
		EmpireI: 0.671,
	}
)

func DetermineHaleAge(realLifeAge int, empire EmpireID) (float64, error) {
	if realLifeAge <= 0 {
		return 0, ers.NewErrInvalidParameter(
			"realLifeAge",
			ers.NewErrGT(0),
		)
	} else if realLifeAge > 83 {
		return 0, ers.NewErrInvalidParameter(
			"realLifeAge",
			ers.NewErrLTE(83),
		)
	}

	age := (float64(realLifeAge) * 27.9 * (1.7 - EmpireUniversalConstants[empire])) / 83.5

	return age, nil
}

func DetermineVusAge(realLifeAge int, empire EmpireID) (float64, error) {
	if realLifeAge <= 0 {
		return 0, ers.NewErrInvalidParameter(
			"realLifeAge",
			ers.NewErrGT(0),
		)
	} else if realLifeAge > 83 {
		return 0, ers.NewErrInvalidParameter(
			"realLifeAge",
			ers.NewErrLTE(83),
		)
	}

	age := (float64(realLifeAge) * 29.01 * (1.7 - EmpireUniversalConstants[empire])) / 83.5

	return age, nil
}
