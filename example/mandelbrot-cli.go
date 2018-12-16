package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	fractal ".."
)

func getAutoZoomMethod(zoomConfig string) func(fractal.Plane) fractal.Plane {
	zoom := strings.Split(zoomConfig, ":")
	zoomType := zoom[0]
	if zoomType == "raster" {
		division, err := strconv.Atoi(zoom[1])
		if err != nil {
			panic(err)
		}
		return func(p fractal.Plane) fractal.Plane {
			return p.RasterAutoZoom(division)
		}
	}
	panic(fmt.Sprintf("unknown zoom type: %s\n", zoomType))
}

func writeFile(num int, outputfile string, image *image.RGBA) {
	f, err := os.Create(fmt.Sprintf("%03d_%s", num, outputfile))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, image)
}

func forQueryParam(r *http.Request, param string, f func(value float64)) {
	values, ok := r.URL.Query()[param]
	if ok {
		fval, err := strconv.ParseFloat(values[0], 64)
		if err == nil {
			f(fval)
		}
	}
}

func main() {
	var xstart float64
	var xend float64
	var ystart float64
	var yend float64
	var iterations int
	var width int
	var height int
	var outputfile string
	var port int
	var zoom int
	var zoomConfig string

	flag.Float64Var(&xstart, "xstart", -2.0, "xstart")
	flag.Float64Var(&xend, "xend", 1.2, "xend")
	flag.Float64Var(&ystart, "ystart", -1.2, "ystart")
	flag.Float64Var(&yend, "yend", 1.2, "yend")
	flag.IntVar(&iterations, "iterations", 100, "iterations")
	flag.IntVar(&width, "width", 400, "width")
	flag.IntVar(&height, "height", 300, "height")
	flag.StringVar(&outputfile, "outputfile", "mandelbrot.png", "outputfile")
	flag.IntVar(&port, "port", 8080, "http port")
	flag.IntVar(&zoom, "zoom", 0, "zoom")
	flag.StringVar(&zoomConfig, "zoom-config", "raster:2", "zoom type [raster:N] default is raster:2")

	flag.Parse()

	m := fractal.ComplexSet{
		fractal.Range{Start: xstart, End: xend},
		fractal.Range{Start: ystart, End: yend},
		fractal.Mandelbrot,
	}

	p := fractal.NewPlane(m, width, height, iterations)
	writeFile(0, outputfile, p.Image())

	autoZoom := getAutoZoomMethod(zoomConfig)
	for z := 0; z < zoom; z++ {
		p = autoZoom(p)
		writeFile(z+1, outputfile, p.Image())
	}
}
