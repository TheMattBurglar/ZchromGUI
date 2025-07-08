package main

import (
	"ZtestAssisted/ztestlogic"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Z Chromosome Simulation")

	// Input fields
	adamEntry := widget.NewEntry()
	adamEntry.SetText("1000")
	eveEntry := widget.NewEntry()
	eveEntry.SetText("1000")
	lilithEntry := widget.NewEntry()
	lilithEntry.SetText("1")
	dianaEntry := widget.NewEntry()
	dianaEntry.SetText("0")

	yEggsCheck := widget.NewCheck("", nil)

	eveBirth := widget.NewEntry()
	eveBirth.SetText("1.5")
	lilithBirth := widget.NewEntry()
	lilithBirth.SetText("1.5")
	dianaBirth := widget.NewEntry()
	dianaBirth.SetText("1.5")

	maxPopEntry := widget.NewEntry()
	maxPopEntry.SetText("50000")
	generationsEntry := widget.NewEntry()
	generationsEntry.SetText("37")
	timelinesEntry := widget.NewEntry()
	timelinesEntry.SetText("100")

	resultLabel := widget.NewLabel("Results will appear here.")

	runButton := widget.NewButton("Run Simulation", func() {
		// Parse input values
		pop := [4]int{}
		birth := [3]float64{}
		pop[0], _ = strconv.Atoi(adamEntry.Text)
		pop[1], _ = strconv.Atoi(eveEntry.Text)
		pop[2], _ = strconv.Atoi(lilithEntry.Text)
		pop[3], _ = strconv.Atoi(dianaEntry.Text)
		birth[0], _ = strconv.ParseFloat(eveBirth.Text, 64)
		birth[1], _ = strconv.ParseFloat(lilithBirth.Text, 64)
		birth[2], _ = strconv.ParseFloat(dianaBirth.Text, 64)
		maxPop, _ := strconv.Atoi(maxPopEntry.Text)
		gens, _ := strconv.Atoi(generationsEntry.Text)
		timelines, _ := strconv.Atoi(timelinesEntry.Text)

		success := 0

		// Reset global counters if needed
		ztestlogic.TotalExtinction = 0
		ztestlogic.Zextinction = 0
		ztestlogic.MaleExtinction = 0
		ztestlogic.FemExtinction = 0
		ztestlogic.LastGen = 0
		ztestlogic.MaxPopReached = 0
		ztestlogic.PopCapGen = 0

		viableY := "N"
		if yEggsCheck.Checked {
			viableY = "Y"
		}

		for i := 0; i < timelines; i++ {
			if ztestlogic.GenTryFail(pop, birth, viableY, maxPop, gens) {
				success++
			}
		}

		// Build the summary string
		summary := ""
		summary += strconv.Itoa(success) + " out of " + strconv.Itoa(timelines) + " timelines still had the Z chromosome by the end.\n"
		if ztestlogic.TotalExtinction > 0 {
			summary += strconv.Itoa(ztestlogic.TotalExtinction) + " failed because EVERYONE died out.\n"
		}
		if ztestlogic.Zextinction > 0 {
			summary += strconv.Itoa(ztestlogic.Zextinction) + " failed because Lilith and Diana died out. There were still Women, but no more Z chromosomes.\n"
		}
		if ztestlogic.MaleExtinction > 0 {
			summary += strconv.Itoa(ztestlogic.MaleExtinction) + " failed because men died out. Usually because total population got too small.\n"
		}
		if ztestlogic.FemExtinction > 0 {
			summary += strconv.Itoa(ztestlogic.FemExtinction) + " failed because women died out completely. Usually because total population got too small.\n"
		}
		if ztestlogic.LastGen > 0 {
			summary += "If they ended without either men or a Z chromosome, they did so within " + strconv.Itoa(ztestlogic.LastGen) + " generations.\n"
		}
		if ztestlogic.MaxPopReached > 0 {
			summary += strconv.Itoa(ztestlogic.MaxPopReached) + " were cut off early because they reached a population size of " + strconv.Itoa(maxPop) + "\n"
			summary += "They hit that population cap at or below " + strconv.Itoa(ztestlogic.PopCapGen) + " generations.\n"
		}

		ZisThere, finalPop := ztestlogic.GenTryFailWithPop(pop, birth, viableY, maxPop, gens)

		isMarker := finalPop[0] == 0 && finalPop[1] == 0 && finalPop[2] == 0 && finalPop[3] > 0
		total := finalPop[0] + finalPop[1] + finalPop[2] + finalPop[3]

		if !isMarker && total > 0 {
			summary += "\nExample timeline final population totals:\n"
			summary += "Adam: " + strconv.Itoa(finalPop[0]) + "\n"
			summary += "Eve: " + strconv.Itoa(finalPop[1]) + "\n"
			summary += "Lilith: " + strconv.Itoa(finalPop[2]) + "\n"
			summary += "Diana: " + strconv.Itoa(finalPop[3]) + "\n"
			summary += "\nPercentages:\n"
			summary += "Adam: " + strconv.FormatFloat(100*float64(finalPop[0])/float64(total), 'f', 2, 64) + "%\n"
			summary += "Eve: " + strconv.FormatFloat(100*float64(finalPop[1])/float64(total), 'f', 2, 64) + "%\n"
			summary += "Lilith: " + strconv.FormatFloat(100*float64(finalPop[2])/float64(total), 'f', 2, 64) + "%\n"
			summary += "Diana: " + strconv.FormatFloat(100*float64(finalPop[3])/float64(total), 'f', 2, 64) + "%\n"
		}
		if ZisThere && isMarker {
			summary += "\nExample timeline ended successfully by reaching the population cap.\n"
		}

		if isMarker && !ZisThere {
			summary += "\nExample timeline ended with a marker indicating Men, Women, or the Z chromosome died out.\n"
		}

		resultLabel.SetText(summary)
	})

	scroll := container.NewVScroll(resultLabel)

	w.SetContent(
		container.NewBorder(
			container.NewVBox(
				widget.NewLabel("Z Chromosome Simulation"),
				container.NewGridWithColumns(2,
					widget.NewLabel("Adam:"), adamEntry,
					widget.NewLabel("Eve:"), eveEntry,
					widget.NewLabel("Lilith:"), lilithEntry,
					widget.NewLabel("Diana:"), dianaEntry,
					widget.NewLabel("Y chromosome eggs viable?"), yEggsCheck, // empty label for alignment
					widget.NewLabel("Eve Birth Rate:"), eveBirth,
					widget.NewLabel("Lilith Birth Rate:"), lilithBirth,
					widget.NewLabel("Diana Birth Rate:"), dianaBirth,
					widget.NewLabel("Max Population:"), maxPopEntry,
					widget.NewLabel("Generations:"), generationsEntry,
					widget.NewLabel("Timelines:"), timelinesEntry,
				),

				runButton,
			), // top
			nil,    // left
			nil,    // right
			nil,    // bottom
			scroll, // center (this will expand!)
		),
	)

	w.ShowAndRun()
}
