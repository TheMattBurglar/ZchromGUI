package main

import (
	"image/color"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// LineChart is a simple widget to display multiple lines of data
type LineChart struct {
	widget.BaseWidget
	Data           [][]float64 // Each inner slice is a line series
	Colors         []color.Color
	MaxY           float64              // Optional manual max Y, if 0 it auto-scales
	MaxX           float64              // Optional manual max X (expected number of points), if 0 it auto-scales
	MinimumSize    fyne.Size            // Minimum size for the chart
	ValueFormatter func(float64) string // Optional formatter for Y-axis values
	XAxisLabel     string               // Label for the X-axis
}

// NewLineChart creates a new LineChart widget
func NewLineChart() *LineChart {
	chart := &LineChart{
		MinimumSize: fyne.NewSize(200, 200), // Increased default height for labels
	}
	chart.ExtendBaseWidget(chart)
	return chart
}

// SetData updates the chart data and refreshes the widget
func (c *LineChart) SetData(data [][]float64) {
	c.Data = data
	c.Refresh()
}

// CreateRenderer implements fyne.Widget
func (c *LineChart) CreateRenderer() fyne.WidgetRenderer {
	return &chartRenderer{chart: c}
}

type chartRenderer struct {
	chart *LineChart
	lines []fyne.CanvasObject
}

func (r *chartRenderer) MinSize() fyne.Size {
	return r.chart.MinimumSize
}

func (r *chartRenderer) Layout(size fyne.Size) {
	// Logic moved to Refresh to handle dynamic line creation/updating
}

func (r *chartRenderer) Refresh() {
	// Clear old objects
	r.lines = nil

	size := r.chart.Size()
	width := float64(size.Width)
	height := float64(size.Height)

	// Margins for labels
	leftMargin := 40.0
	bottomMargin := 40.0 // Increased for X-axis label
	chartWidth := width - leftMargin
	chartHeight := height - bottomMargin

	// Draw background
	bg := canvas.NewRectangle(theme.BackgroundColor())
	bg.Resize(size)
	r.lines = append(r.lines, bg)

	// Draw Axes
	axisColor := theme.ForegroundColor()

	// Y Axis
	yAxis := canvas.NewLine(axisColor)
	yAxis.Position1 = fyne.NewPos(float32(leftMargin), 0)
	yAxis.Position2 = fyne.NewPos(float32(leftMargin), float32(chartHeight))
	yAxis.StrokeWidth = 1
	r.lines = append(r.lines, yAxis)

	// X Axis
	xAxis := canvas.NewLine(axisColor)
	xAxis.Position1 = fyne.NewPos(float32(leftMargin), float32(chartHeight))
	xAxis.Position2 = fyne.NewPos(float32(width), float32(chartHeight))
	xAxis.StrokeWidth = 1
	r.lines = append(r.lines, xAxis)

	if len(r.chart.Data) == 0 && r.chart.MaxX == 0 {
		return
	}

	// Determine ranges
	maxX := 0.0
	maxY := r.chart.MaxY

	if r.chart.MaxX > 0 {
		maxX = r.chart.MaxX
	}

	for _, series := range r.chart.Data {
		if float64(len(series)) > maxX {
			maxX = float64(len(series))
		}
		if r.chart.MaxY == 0 {
			for _, val := range series {
				if val > maxY {
					maxY = val
				}
			}
		}
	}

	if maxX < 2 || maxY <= 0 {
		return
	}

	// Draw Labels
	// Y Max
	yMaxStr := strconv.FormatFloat(maxY, 'f', 0, 64)
	if r.chart.ValueFormatter != nil {
		yMaxStr = r.chart.ValueFormatter(maxY)
	}
	yMaxText := canvas.NewText(yMaxStr, theme.ForegroundColor())
	yMaxText.Alignment = fyne.TextAlignTrailing
	yMaxText.TextSize = 10
	yMaxText.Move(fyne.NewPos(float32(leftMargin)-5, 0)) // Slightly left of axis
	r.lines = append(r.lines, yMaxText)

	// Y Zero
	yZeroStr := "0"
	if r.chart.ValueFormatter != nil {
		yZeroStr = r.chart.ValueFormatter(0)
	}
	yZeroText := canvas.NewText(yZeroStr, theme.ForegroundColor())
	yZeroText.Alignment = fyne.TextAlignTrailing
	yZeroText.TextSize = 10
	yZeroText.Move(fyne.NewPos(float32(leftMargin)-5, float32(chartHeight)-10))
	r.lines = append(r.lines, yZeroText)

	// X Max (Generations)
	xMaxText := canvas.NewText(strconv.Itoa(int(maxX)-1), theme.ForegroundColor())
	xMaxText.Alignment = fyne.TextAlignCenter
	xMaxText.TextSize = 10
	xMaxText.Move(fyne.NewPos(float32(width)-20, float32(chartHeight)+5))
	r.lines = append(r.lines, xMaxText)

	// X Zero
	xZeroText := canvas.NewText("0", theme.ForegroundColor())
	xZeroText.Alignment = fyne.TextAlignCenter
	xZeroText.TextSize = 10
	xZeroText.Move(fyne.NewPos(float32(leftMargin), float32(chartHeight)+5))
	r.lines = append(r.lines, xZeroText)

	// X Axis Label
	if r.chart.XAxisLabel != "" {
		xLabel := canvas.NewText(r.chart.XAxisLabel, theme.ForegroundColor())
		xLabel.Alignment = fyne.TextAlignCenter
		xLabel.TextSize = 12
		xLabel.TextStyle = fyne.TextStyle{Bold: true}
		// Center horizontally in the chart area
		xLabel.Move(fyne.NewPos(float32(leftMargin+chartWidth/2), float32(chartHeight)+20))
		r.lines = append(r.lines, xLabel)
	}

	// Draw lines
	for i, series := range r.chart.Data {
		if len(series) < 2 {
			continue
		}

		lineColor := theme.PrimaryColor()
		if len(r.chart.Colors) > i {
			lineColor = r.chart.Colors[i]
		}

		xStep := chartWidth / float64(maxX-1)
		yScale := chartHeight / maxY

		for j := 0; j < len(series)-1; j++ {
			x1 := leftMargin + float64(j)*xStep
			y1 := chartHeight - (series[j] * yScale)
			x2 := leftMargin + float64(j+1)*xStep
			y2 := chartHeight - (series[j+1] * yScale)

			segment := canvas.NewLine(lineColor)
			segment.Position1 = fyne.NewPos(float32(x1), float32(y1))
			segment.Position2 = fyne.NewPos(float32(x2), float32(y2))
			segment.StrokeWidth = 1.5
			r.lines = append(r.lines, segment)
		}
	}
}

func (r *chartRenderer) Objects() []fyne.CanvasObject {
	return r.lines
}

func (r *chartRenderer) Destroy() {
}

// Helper to downsample data if it's too large for the screen width
func downsample(data []float64, maxPoints int) []float64 {
	if len(data) <= maxPoints {
		return data
	}
	step := float64(len(data)) / float64(maxPoints)
	result := make([]float64, maxPoints)
	for i := 0; i < maxPoints; i++ {
		index := int(math.Round(float64(i) * step))
		if index >= len(data) {
			index = len(data) - 1
		}
		result[i] = data[index]
	}
	return result
}
