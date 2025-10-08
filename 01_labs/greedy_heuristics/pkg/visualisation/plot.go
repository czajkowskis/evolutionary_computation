package visualisation

import (
	"fmt"
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"github.com/czajkowskis/evolutionary_computation/01_labs/greedy_heuristics/pkg/data"
)

// PlotSolution draws a TSP path with (0,0) axes, correct arrowheads, and square scaling.
func PlotSolution(nodes []data.Node, path []int, filename string,
	xMin, xMax, yMin, yMax float64) error {

	p := plot.New()
	p.Title.Text = "TSP Solution"

	// --- Prepare path points ---
	pathPoints := make(plotter.XYs, len(path)+1)
	for i, idx := range path {
		pathPoints[i].X = float64(nodes[idx].X)
		pathPoints[i].Y = float64(nodes[idx].Y)
	}
	if len(path) > 0 {
		pathPoints[len(path)] = pathPoints[0]
	}

	// --- Prepare all node points ---
	allPoints := make(plotter.XYs, len(nodes))
	for i, n := range nodes {
		allPoints[i].X = float64(n.X)
		allPoints[i].Y = float64(n.Y)
	}

	// --- Path line ---
	line, _ := plotter.NewLine(pathPoints)
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = vg.Points(1.5)

	// --- Scatter plots ---
	allScatter, _ := plotter.NewScatter(allPoints)
	allScatter.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	allScatter.GlyphStyle.Radius = vg.Points(3)

	pathScatter, _ := plotter.NewScatter(pathPoints[:len(path)])
	pathScatter.GlyphStyle.Color = color.RGBA{G: 255, A: 255}
	pathScatter.GlyphStyle.Radius = vg.Points(4)

	p.Add(line, allScatter, pathScatter)

	// --- Determine axis bounds ---
	if xMin > 0 {
		xMin = 0
	}
	if yMin > 0 {
		yMin = 0
	}

	// enforce a square coordinate area (equal scaling)
	p.X.Min, p.Y.Min = 0, 0
	p.X.Max, p.Y.Max = xMax, yMax

	// --- Draw X and Y axes (through origin) ---
	xAxisPts := plotter.XYs{{X: xMin, Y: 0}, {X: xMax, Y: 0}}
	yAxisPts := plotter.XYs{{X: 0, Y: yMin}, {X: 0, Y: yMax}}

	xAxis, _ := plotter.NewLine(xAxisPts)
	yAxis, _ := plotter.NewLine(yAxisPts)
	xAxis.Color, yAxis.Color = color.Black, color.Black
	xAxis.Width, yAxis.Width = vg.Points(1), vg.Points(1)

	p.Add(xAxis, yAxis)

	// --- Save the plot ---
	if err := p.Save(6*vg.Inch, 6*vg.Inch, fmt.Sprintf("%s.png", filename)); err != nil {
		return err
	}

	return nil
}
