package config

// DefaultPlotBounds defines default plot dimensions for TSP instances
var DefaultPlotBounds = PlotBounds{
	XMin: 0,
	XMax: 4000,
	YMin: 0,
	YMax: 2000,
}

// PlotBounds defines the bounds for plotting solutions
type PlotBounds struct {
	XMin, XMax, YMin, YMax float64
}

// InstancePaths holds paths to instance files
type InstancePaths struct {
	TSPA string
	TSPB string
}

// DefaultInstancePaths returns default paths relative to lab directories
func DefaultInstancePaths() InstancePaths {
	return InstancePaths{
		TSPA: "../../instances/TSPA.csv",
		TSPB: "../../instances/TSPB.csv",
	}
}

// GetInstancePath returns the path for a given instance name
func GetInstancePath(instanceName string) string {
	if instanceName == "A" {
		return DefaultInstancePaths().TSPA
	}
	return DefaultInstancePaths().TSPB
}

