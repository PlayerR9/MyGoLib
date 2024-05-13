package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	pkg "github.com/PlayerR9/MyGoLib/pkg"

	pnl "github.com/PlayerR9/MyGoLib/ComplexData/ConsolePanel"
	doc "github.com/PlayerR9/MyGoLib/CustomData/Document"
)

var Console *pnl.ConsolePanel = MakeConsolePanel()

func main() {
	command, err := Console.ParseArguments([]string{"help"}) // []string{"age_converter", "--age", "20"})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = command.Execute()
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Command executed successfully")
	}
}

func CommandAgeConverter(args map[string]any) error {
	// Get flags
	age := args["--age"].(int)

	var empires []pkg.EmpireID

	empire, ok := args["--empire"]
	if !ok {
		empires = []pkg.EmpireID{
			pkg.EmpireA,
			pkg.EmpireB,
			pkg.EmpireC,
			pkg.EmpireD,
			pkg.EmpireE,
			pkg.EmpireF,
			pkg.EmpireG,
			pkg.EmpireH,
			pkg.EmpireI,
		}
	} else {
		empires = []pkg.EmpireID{empire.(pkg.EmpireID)}
	}

	var gender pkg.Gender

	isMale, ok := args["--gender"].(bool)
	if !ok {
		gender = pkg.GenderUnknown
	} else if isMale {
		gender = pkg.GenderVus
	} else {
		gender = pkg.GenderHale
	}

	switch gender {
	case pkg.GenderHale:
		for _, empire := range empires {
			newAge, err := pkg.DetermineHaleAge(age, empire)
			if err != nil {
				return err
			}

			fmt.Printf("The age in the empire %s is %.0f\n", empire.String(), newAge)
		}
	case pkg.GenderVus:
		for _, empire := range empires {
			newAge, err := pkg.DetermineVusAge(age, empire)
			if err != nil {
				return err
			}

			fmt.Printf("The age in the empire %s is %.0f\n", empire.String(), newAge)
		}
	case pkg.GenderUnknown:
		fmt.Println("For female characters:")

		for _, empire := range empires {
			newAge, err := pkg.DetermineHaleAge(age, empire)
			if err != nil {
				return err
			}

			fmt.Printf("The age in the empire %s is %.0f\n", empire.String(), newAge)
		}

		fmt.Println("For male characters:")

		for _, empire := range empires {
			newAge, err := pkg.DetermineVusAge(age, empire)
			if err != nil {
				return err
			}

			fmt.Printf("The age in the empire %s is %.0f\n", empire.String(), newAge)
		}
	}

	return nil
}

func MakeConsolePanel() *pnl.ConsolePanel {
	console := pnl.NewConsolePanel(
		os.Args[0],
		doc.NewDocument(
			"This command converts the age of a person in real life to",
			"the age of the person in the empire of the user's choice.",
		).AddLine(
			"This tool is useful for the conworld project \"2S\" as it",
			"allows to determine, lore-wise, how old a person is in a",
			"specific empire.",
		),
	).AddCommand(
		"age_converter",
		pnl.NewCommandInfo(
			doc.NewDocument(
				"This command converts the age of a person in real life to",
				"the age of the person in the empire of the user's choice.",
			).AddLine(
				"This tool is useful for the conworld project \"2S\" as it",
				"allows to determine, lore-wise, how old a person is in a",
				"specific empire.",
			),
			CommandAgeConverter,
		).AddFlag(
			"--age",
			pnl.NewFlagInfo(
				true,
				nil,
				pnl.NewArgument(
					"age",
					func(s string) (any, error) {
						age, err := strconv.Atoi(s)
						if err != nil {
							return nil, fmt.Errorf("age must be an integer")
						}

						if age <= 0 {
							return nil, fmt.Errorf("age must be greater than 0")
						} else if age > 83 {
							return nil, fmt.Errorf("age must be less than or equal to 83")
						}

						return age, nil
					},
				),
			).SetDescription(
				doc.NewDocument(
					"The age of the person in real life.",
				).AddLine(
					"It must be an integer greater than 0 and less than or equal to 83.",
				),
			),
		).AddFlag(
			"--gender",
			pnl.NewFlagInfo(
				false,
				nil,
				pnl.NewArgument(
					"gender",
					func(s string) (any, error) {
						choice := strings.ToLower(s)

						if choice == "m" || choice == "male" {
							return true, nil
						} else if choice == "f" || choice == "female" {
							return false, nil
						}

						return nil, fmt.Errorf("gender value (%s) is not valid", s)
					},
				),
			).SetDescription(
				doc.NewDocument(
					"The gender is used to apply the correct aging factor to the age.",
				).AddLine(
					"It must be either \"M\" or \"Male\" for male characters and",
					"\"F\" or \"Female\" for female characters.",
				).AddLine(
					"The choice is not case-sensitive and, when omitted, both genders",
					"will be shown.",
				),
			),
		).AddFlag(
			"--empire",
			pnl.NewFlagInfo(
				false,
				nil,
				pnl.NewArgument(
					"empire",
					func(s string) (any, error) {
						empire, err := pkg.StringToEmpireID(s)
						if err != nil {
							return nil, err
						}

						return empire, nil
					},
				),
			).SetDescription(
				doc.NewDocument(
					"The empire to convert the age to.",
				).AddLine(
					"It must be one of the following: A, B, C, D, E, F, G, H, and I, and",
					"it is not case-sensitive.",
				).AddLine(
					"When omitted, the age will be converted to all empires.",
				),
			),
		),
	)

	return console
}
