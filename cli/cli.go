package cli

import (
	"fmt"
	"os"

	"../batten"
	"github.com/mgutz/ansi"
	"github.com/olekukonko/tablewriter"
)

var (
	lime       = ansi.ColorCode("green+h:black")
	red        = ansi.ColorCode("red")
	green      = ansi.ColorCode("green")
	redonwhite = ansi.ColorCode("red:white")
	yellow     = ansi.ColorCode("yellow")
	reset      = ansi.ColorCode("reset")
)

var resultsError = red + "FAILED (error)" + reset
var resultsFailed = red + "FAILED" + reset
var resultsOK = lime + "PASSED" + reset

//
// FormatResultsForConsole formats the `CheckResults` for
// console display.
//
func FormatResultsForConsole(idx int, results *batten.CheckResults) {

	checkdefinition := results.CheckDefinition

	fmt.Printf("[%d/%d] ", idx+1, len(batten.Checks))

	if results.Error != nil {
		fmt.Printf("%s [%s] %s\n", resultsError, checkdefinition.Identifier(), checkdefinition.Name())
		fmt.Println("\t There was an error executing the check:", results.Error)
	} else {
		if results.Success {
			fmt.Printf("%s [%s] %s\n", resultsOK, checkdefinition.Identifier(), checkdefinition.Name())
		} else {
			fmt.Printf("%s [%s] %s\n", resultsFailed, checkdefinition.Identifier(), checkdefinition.Name())
			table := tablewriter.NewWriter(os.Stdout)
			table.SetBorder(false)
			table.SetColWidth(75)
			table.Append([]string{
				ansi.LightWhite + "Description" + reset,
				checkdefinition.Description(),
			})
			table.Append([]string{
				ansi.LightWhite + "Remediation" + reset,
				checkdefinition.Remediation(),
			})

			table.Render()
		}
	}
}
