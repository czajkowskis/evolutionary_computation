package visualisation

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"path/filepath"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"github.com/czajkowskis/evolutionary_computation/02_labs/greedy_regret_heuristics/pkg/data"
)

// PlotSolution draws a TSP path with (0,0) axes, correct arrowheads, square scaling and cost-scaled node sizes and a blue gradient.
func PlotSolution(nodes []data.Node, path []int, title string, filename string,
	xMin, xMax, yMin, yMax float64) error {

	p := plot.New()

	// --- Add title inside the plot ---

	xCenter := xMin + (xMax-xMin)/2
	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: xCenter, Y: yMax},
		},
		Labels: []string{title},
	})
	if err != nil {
		return err
	}
	labels.Offset.Y = vg.Points(25)
	if len(labels.TextStyle) > 0 {
		labels.TextStyle[0].Font.Size = vg.Points(14)
		labels.TextStyle[0].XAlign = draw.XCenter
		labels.TextStyle[0].YAlign = draw.YTop
	}
	p.Add(labels)

	// Move axes to origin
	p.X.LineStyle.Width = 0
	p.Y.LineStyle.Width = 0

	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

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

	// ---- Cost range  ----
	minCost, maxCost := nodes[0].Cost, nodes[0].Cost
	for _, n := range nodes {
		if n.Cost < minCost {
			minCost = n.Cost
		}
		if n.Cost > maxCost {
			maxCost = n.Cost
		}
	}

	// --- Path line ---
	line, _ := plotter.NewLine(pathPoints)
	line.Color = color.RGBA{R: 255, A: 255}
	line.Width = vg.Points(1.5)

	// ---- Scatter with size and color by cost ----
	scatter, _ := plotter.NewScatter(allPoints)
	scatter.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		r := scaleCostToRadius(nodes[i].Cost, minCost, maxCost, vg.Points(3.5), vg.Points(6.5))
		c := costToBlue(nodes[i].Cost, minCost, maxCost) // light -> dark blue gradient
		return draw.GlyphStyle{
			Color:  c,
			Radius: r,
			Shape:  draw.CircleGlyph{},
		}
	}

	// ---- Highlight path nodes as small black crosses so they stand out ----
	var pathScatter *plotter.Scatter
	if len(path) > 0 {
		pathScatter, _ = plotter.NewScatter(pathPoints[:len(path)])
		pathScatter.GlyphStyle.Color = color.Black
		pathScatter.GlyphStyle.Radius = vg.Points(3.0)
		pathScatter.GlyphStyle.Shape = draw.CrossGlyph{}
	}

	p.Add(line, scatter)
	if pathScatter != nil {
		p.Add(pathScatter)
	}
	p.Add(plotter.NewGrid())

	// --- Determine axis bounds ---
	if xMin > 0 {
		xMin = 0
	}
	if yMin > 0 {
		yMin = 0
	}

	// enforce a square coordinate area (equal scaling)
	xRange := xMax - xMin
	yRange := yMax - yMin
	maxRange := math.Max(xRange, yRange)
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

	// --- Add arrowheads at the positive ends ---
	arrowSize := maxRange * 0.02 // 2 % of the larger axis range

	// --- X-axis arrow  ---
	ax1Pts := plotter.XYs{{X: xMax - arrowSize, Y: 0.25 * arrowSize}, {X: xMax, Y: 0}}
	ax2Pts := plotter.XYs{{X: xMax - arrowSize, Y: -0.25 * arrowSize}, {X: xMax, Y: 0}}
	ax1, _ := plotter.NewLine(ax1Pts)
	ax2, _ := plotter.NewLine(ax2Pts)
	ax1.Color, ax2.Color = color.Black, color.Black
	p.Add(ax1, ax2)

	// --- Y-axis arrow  ---
	// compensate for aspect ratio so the arrow looks the same visually
	aspect := (xMax - xMin) / (yMax - yMin)
	lengthRatio := 0.75 // slight shortening to visually match X arrow
	ay1Pts := plotter.XYs{{X: -0.25 * arrowSize * aspect, Y: yMax - arrowSize*lengthRatio}, {X: 0, Y: yMax}}
	ay2Pts := plotter.XYs{{X: 0.25 * arrowSize * aspect, Y: yMax - arrowSize*lengthRatio}, {X: 0, Y: yMax}}
	ay1, _ := plotter.NewLine(ay1Pts)
	ay2, _ := plotter.NewLine(ay2Pts)
	ay1.Color, ay2.Color = color.Black, color.Black
	p.Add(ay1, ay2)

	// --- Save the plot ---
	plotDir := "output_plots"
	if err := os.MkdirAll(plotDir, 0755); err != nil {
		return err
	}
	filePath := filepath.Join(plotDir, fmt.Sprintf("%s.png", filename))
	if err := p.Save(8*vg.Inch, 8*vg.Inch, filePath); err != nil {
		return err
	}

	return nil
}

// Size scaling: map int cost -> radius in [minR, maxR].
func scaleCostToRadius(cost, minCost, maxCost int, minR, maxR vg.Length) vg.Length {
	if maxCost == minCost {
		return (minR + maxR) / 2
	}
	n := float64(cost-minCost) / float64(maxCost-minCost)
	return minR + vg.Length(n)*(maxR-minR)
}

// Single-hue blue gradient (light -> dark) low: #c6dbef, high: #084594.
func costToBlue(cost, minCost, maxCost int) color.RGBA {
	low := color.RGBA{R: 0xC6, G: 0xDB, B: 0xEF, A: 0xFF}  // light blue
	high := color.RGBA{R: 0x08, G: 0x45, B: 0x94, A: 0xFF} // dark blue
	if maxCost == minCost {
		return low
	}
	n := float64(cost-minCost) / float64(maxCost-minCost) // 0..1
	return lerpRGBA(low, high, n)
}

func lerpRGBA(a, b color.RGBA, t float64) color.RGBA {
	clamp := func(x float64) uint8 {
		if x < 0 {
			return 0
		}
		if x > 255 {
			return 255
		}
		return uint8(x + 0.5)
	}
	return color.RGBA{
		R: clamp(float64(a.R) + (float64(b.R)-float64(a.R))*t),
		G: clamp(float64(a.G) + (float64(b.G)-float64(a.G))*t),
		B: clamp(float64(a.B) + (float64(b.B)-float64(a.B))*t),
		A: 255,
	}
}
