package main

import (
	"ZtestAssisted/ztestlogic"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()

	// Apply custom theme
	a.Settings().SetTheme(&myTheme{})

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

	yEggsCheck := widget.NewCheck("Y chromosome eggs viable?", nil)

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
	resultLabel.Wrapping = fyne.TextWrapWord

	// Charts
	menChart := NewLineChart()
	menChart.Colors = []color.Color{color.RGBA{R: 0x03, G: 0xda, B: 0xc6, A: 0xff}} // Teal
	menChart.ValueFormatter = func(v float64) string {
		return strconv.FormatFloat(v, 'f', 1, 64) + "%"
	}
	menChart.XAxisLabel = "Generations"

	totalChart := NewLineChart()
	totalChart.Colors = []color.Color{color.RGBA{R: 0xbb, G: 0x86, B: 0xfc, A: 0xff}} // Purple
	totalChart.XAxisLabel = "Generations"

	zChart := NewLineChart()
	zChart.Colors = []color.Color{color.RGBA{R: 0xcf, G: 0x66, B: 0x79, A: 0xff}} // Red/Pink
	zChart.ValueFormatter = func(v float64) string {
		return strconv.FormatFloat(v, 'f', 1, 64) + "%"
	}
	zChart.XAxisLabel = "Generations"

	// Wrap charts in cards or containers with labels
	menChartContainer := container.NewPadded(container.NewVBox(widget.NewLabel("Men Population (%)"), container.NewPadded(menChart)))
	totalChartContainer := container.NewPadded(container.NewVBox(widget.NewLabel("Total Population"), container.NewPadded(totalChart)))
	zChartContainer := container.NewPadded(container.NewVBox(widget.NewLabel("Z-Carriers (%) (Lilith+Diana)"), container.NewPadded(zChart)))

	// Set min size for charts
	menChart.MinimumSize = fyne.NewSize(0, 150)
	totalChart.MinimumSize = fyne.NewSize(0, 150)
	zChart.MinimumSize = fyne.NewSize(0, 150)

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	var runButton *widget.Button
	runButton = widget.NewButton("Run Simulation", func() {
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

		viableY := "N"
		if yEggsCheck.Checked {
			viableY = "Y"
		}

		// Disable button and show progress
		runButton.Disable()
		progressBar.SetValue(0)
		progressBar.Show()
		resultLabel.SetText("Running simulation...")

		// Clear old charts
		menChart.SetData(nil)
		totalChart.SetData(nil)
		zChart.SetData(nil)

		// Set MaxX for charts to ensure X-axis spans the full requested generation count
		// +1 because generations start at 0
		menChart.MaxX = float64(gens + 1)
		totalChart.MaxX = float64(gens + 1)
		zChart.MaxX = float64(gens + 1)

		// Run in a goroutine to keep UI responsive
		go func() {
			success := 0
			stats := &ztestlogic.SimulationStats{}

			// Data for charts
			var menHistory [][]float64
			var totalHistory [][]float64
			var zHistory [][]float64

			// We need to store all histories first to separate the "aura" from the "example"
			// We will take at most N random samples for the aura to avoid performance issues if timelines is huge
			// and always ensure the LAST one is the example.
			type TimelineResult struct {
				Men   []float64
				Total []float64
				Z     []float64
			}
			var allResults []TimelineResult

			for i := 0; i < timelines; i++ {
				// Use the new History function
				succeeded, hist := ztestlogic.GenTryFailHistory(pop, birth, viableY, maxPop, gens, stats)
				if succeeded {
					success++
				}

				// Prepare data for this timeline
				var menLine []float64
				var totalLine []float64
				var zLine []float64

				for _, gen := range hist {
					men := float64(gen[0])
					eve := float64(gen[1])
					lilith := float64(gen[2])
					diana := float64(gen[3])

					total := men + eve + lilith + diana

					menPct := 0.0
					zPct := 0.0
					if total > 0 {
						menPct = (men / total) * 100
						zPct = ((lilith + diana) / total) * 100
					}

					menLine = append(menLine, menPct)
					totalLine = append(totalLine, total)
					zLine = append(zLine, zPct)
				}
				allResults = append(allResults, TimelineResult{Men: menLine, Total: totalLine, Z: zLine})

				// Update progress
				currentProgress := float64(i+1) / float64(timelines)
				fyne.Do(func() {
					progressBar.SetValue(currentProgress)
				})
			}

			// Now prepare the chart data arrays
			// We want: [Aura 1, Aura 2, ..., Example]
			// Limit Aura to ~50 lines for performance, plus the last one as Example

			// Define Colors
			// Men: Teal
			menBase := color.RGBA{R: 0x03, G: 0xda, B: 0xc6, A: 0xff}
			menAura := color.RGBA{R: 0x03, G: 0xda, B: 0xc6, A: 0x15} // Low opacity

			// Total: Purple
			totalBase := color.RGBA{R: 0xbb, G: 0x86, B: 0xfc, A: 0xff}
			totalAura := color.RGBA{R: 0xbb, G: 0x86, B: 0xfc, A: 0x15}

			// Z: Red/Pink
			zBase := color.RGBA{R: 0xcf, G: 0x66, B: 0x79, A: 0xff}
			zAura := color.RGBA{R: 0xcf, G: 0x66, B: 0x79, A: 0x15}

			var menColors []color.Color
			var menWidths []float32

			var totalColors []color.Color
			var totalWidths []float32

			var zColors []color.Color
			var zWidths []float32

			// Helper to add a trace
			addTrace := func(res TimelineResult, isExample bool) {
				menHistory = append(menHistory, res.Men)
				totalHistory = append(totalHistory, res.Total)
				zHistory = append(zHistory, res.Z)

				width := float32(1.0)
				if isExample {
					width = 3.0
					menColors = append(menColors, menBase)
					totalColors = append(totalColors, totalBase)
					zColors = append(zColors, zBase)
				} else {
					menColors = append(menColors, menAura)
					totalColors = append(totalColors, totalAura)
					zColors = append(zColors, zAura)
				}
				menWidths = append(menWidths, width)
				totalWidths = append(totalWidths, width)
				zWidths = append(zWidths, width)
			}

			count := len(allResults)
			// Decide which indices to show.
			// We ALWAYS show the last one (index count-1) as the Example.
			// Depending on count, we sample others.

			// Indices to include in aura
			var auraIndices []int
			maxAura := 50
			if count > 1 {
				step := 1
				if count-1 > maxAura {
					step = (count - 1) / maxAura
				}
				for i := 0; i < count-1; i += step {
					auraIndices = append(auraIndices, i)
				}
			}

			// Add Aura traces
			for _, idx := range auraIndices {
				addTrace(allResults[idx], false)
			}

			// Add Example trace (last one)
			if count > 0 {
				addTrace(allResults[count-1], true)
			}

			// Build the summary string
			summary := ""
			summary += strconv.Itoa(success) + " out of " + strconv.Itoa(timelines) + " timelines still had the Z chromosome by the end.\n"
			if stats.TotalExtinction > 0 {
				summary += strconv.Itoa(stats.TotalExtinction) + " failed because EVERYONE died out.\n"
			}
			if stats.Zextinction > 0 {
				summary += strconv.Itoa(stats.Zextinction) + " failed because Lilith and Diana died out. There were still Women, but no more Z chromosomes.\n"
			}
			if stats.MaleExtinction > 0 {
				summary += strconv.Itoa(stats.MaleExtinction) + " failed because men died out. Usually because total population got too small.\n"
			}
			if stats.FemExtinction > 0 {
				summary += strconv.Itoa(stats.FemExtinction) + " failed because women died out completely. Usually because total population got too small.\n"
			}
			if stats.LastGen > 0 {
				summary += "If they ended without either men or a Z chromosome, they did so within " + strconv.Itoa(stats.LastGen) + " generations.\n"
			}
			if stats.MaxPopReached > 0 {
				summary += strconv.Itoa(stats.MaxPopReached) + " were cut off early because they reached a population size of " + strconv.Itoa(maxPop) + "\n"
				summary += "They hit that population cap at or below " + strconv.Itoa(stats.PopCapGen) + " generations.\n"
			}

			ZisThere, finalPop := ztestlogic.GenTryFailWithPop(pop, birth, viableY, maxPop, gens, stats)

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

			// Update UI on main thread
			fyne.Do(func() {
				resultLabel.SetText(summary)

				menChart.Colors = menColors
				menChart.StrokeWidths = menWidths
				menChart.SetData(menHistory)

				totalChart.Colors = totalColors
				totalChart.StrokeWidths = totalWidths
				totalChart.SetData(totalHistory)

				zChart.Colors = zColors
				zChart.StrokeWidths = zWidths
				zChart.SetData(zHistory)

				runButton.Enable()
				progressBar.Hide()
			})
		}()
	})

	// Right side content: Tabs for Summary and Graphs
	// Use a Grid for the right content to ensure full width usage and better resizing behavior
	rightContent := container.NewVScroll(container.NewVBox(
		container.NewPadded(resultLabel), // Pad the label
		widget.NewSeparator(),
		menChartContainer,
		totalChartContainer,
		zChartContainer,
	))

	// Layout Organization
	initialPopCard := widget.NewCard("Initial Population", "", container.NewGridWithColumns(2,
		widget.NewLabel("Adam:"), adamEntry,
		widget.NewLabel("Eve:"), eveEntry,
		widget.NewLabel("Lilith:"), lilithEntry,
		widget.NewLabel("Diana:"), dianaEntry,
	))

	birthRatesCard := widget.NewCard("Birth Rates", "", container.NewGridWithColumns(2,
		widget.NewLabel("Eve:"), eveBirth,
		widget.NewLabel("Lilith:"), lilithBirth,
		widget.NewLabel("Diana:"), dianaBirth,
	))

	settingsCard := widget.NewCard("Simulation Settings", "", container.NewVBox(
		container.NewGridWithColumns(2,
			widget.NewLabel("Max Population:"), maxPopEntry,
			widget.NewLabel("Generations:"), generationsEntry,
			widget.NewLabel("Timelines:"), timelinesEntry,
		),
		yEggsCheck,
	))

	// Left side content: Scrollable inputs + Fixed bottom button
	leftContent := container.NewBorder(
		widget.NewLabelWithStyle("Z Chromosome Simulation", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), // top
		container.NewVBox(progressBar, runButton),                                                             // bottom
		nil, // left
		nil, // right
		container.NewVScroll(container.NewVBox( // center (scrollable)
			initialPopCard,
			birthRatesCard,
			settingsCard,
		)),
	)

	w.SetContent(
		container.NewHSplit(
			leftContent,
			rightContent,
		),
	)

	w.Resize(fyne.NewSize(1000, 700)) // Increased size for charts
	w.ShowAndRun()
}

// Custom Theme
type myTheme struct{}

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 0x12, G: 0x12, B: 0x12, A: 0xff}
	case theme.ColorNameForeground:
		return color.RGBA{R: 0xe0, G: 0xe0, B: 0xe0, A: 0xff}
	case theme.ColorNamePrimary:
		return color.RGBA{R: 0xbb, G: 0x86, B: 0xfc, A: 0xff}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 0x2c, G: 0x2c, B: 0x2c, A: 0xff}
	case theme.ColorNamePlaceHolder:
		return color.RGBA{R: 0x88, G: 0x88, B: 0x88, A: 0xff}
	case theme.ColorNameFocus:
		return color.RGBA{R: 0x03, G: 0xda, B: 0xc6, A: 0xff}
	case theme.ColorNameScrollBar:
		return color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0x00}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
