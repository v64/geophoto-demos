package main

import (
	"flag"
	"fmt"
	"github.com/v64/geophoto"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

func main() {
	flag.Parse()

	inDir := flag.Arg(0)
	if inDir == "" {
		fmt.Fprintf(os.Stderr, "Input directory is required\n")
		os.Exit(1)
	}

	outDir := flag.Arg(1)
	if outDir == "" {
		fmt.Fprintf(os.Stderr, "Output directory is required\n")
		os.Exit(1)
	}

	if !strings.HasSuffix(outDir, "/") {
		outDir += "/"
	}

	_, err := os.Stat(outDir)
	if err == nil {
		fmt.Fprintf(os.Stderr, "Output directory already exists\n")
		os.Exit(1)
	}

	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output directory: %s\n", err.Error())
		os.Exit(1)
	}

	geoDatas := geophoto.DirGeoPhotoData(inDir)

	var timestamps []int
	for timestamp, _ := range geoDatas {
		timestamps = append(timestamps, timestamp)
	}

	sort.Ints(timestamps)
	total := len(timestamps)

	for i, timestamp := range timestamps {
		geo := geoDatas[timestamp]
		fmt.Println("Downloading", i+1, "of", total)
		getStreetViewImage(outDir, fmt.Sprintf("%05d", i+1), getStreetViewUrl(geo))
	}
}

func getStreetViewImage(path, name, url string) {
	fileName := path + name + ".jpg"

	out, err := os.Create(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating file %s: %s\n", fileName, err.Error())
		return
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting URL %s: %s\n", url, err.Error())
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file %s: %s\n", url, err.Error())
		return
	}

	// Avoid hitting google's rate limit
	time.Sleep(time.Second)
}

func getStreetViewUrl(geo geophoto.GeoPhoto) string {
	return fmt.Sprintf("https://maps.googleapis.com/maps/api/streetview?location=%s&size=640x640&fov=120&heading=0&sensor=false", geo.StringDegrees())
}
