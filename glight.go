package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xplshn/a-utils/pkg/ccmd"

	"github.com/blackjack/webcam"
	_ "image/jpeg"
)

// AnalyzeBrightness calculates the average brightness of an image and scales it to a range between 1 and 100.
func analyzeBrightness(img image.Image) float64 {
	bounds := img.Bounds()
	var totalBrightness float64
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			luminance := 0.299*float64(r)/0xFFFF + 0.587*float64(g)/0xFFFF + 0.114*float64(b)/0xFFFF
			totalBrightness += luminance
		}
	}
	averageBrightness := totalBrightness / float64(bounds.Dx()*bounds.Dy())

	// Scale the average brightness to a range between 1 and 100
	scaledBrightness := (averageBrightness * 99.0) + 1.0
	return scaledBrightness
}

// setBrightness adjusts the screen brightness with smooth transition.
func setBrightness(level float64, brightnessFile, maxBrightnessFile string, minBrightness uint8, scaleFactor int) {
	if level == 0 {
		log.Fatalf("Refusing to set brightness to zero: %v", level)
	}
	maxBrightnessData, err := os.ReadFile(maxBrightnessFile)
	if err != nil {
		log.Fatalf("Failed to read max brightness: %v", err)
	}
	maxBrightness := atoi(string(maxBrightnessData))

	currentData, err := os.ReadFile(brightnessFile)
	if err != nil {
		log.Fatalf("Failed to read current brightness: %v", err)
	}
	currentBrightness := atoi(string(currentData))

	newBrightness := int(level) * maxBrightness / 100
	if newBrightness < int(minBrightness) * maxBrightness / 100 {
		newBrightness = int(minBrightness) * maxBrightness / 100
	}

	// Smooth transition with scaling factor
	step := scaleFactor
	if newBrightness < currentBrightness {
		step = -scaleFactor
	}
	for i := currentBrightness; (step > 0 && i < newBrightness) || (step < 0 && i > newBrightness); i += step {
		err = os.WriteFile(brightnessFile, []byte(fmt.Sprintf("%d", i)), 0644)
		if err != nil {
			log.Fatalf("Failed to set brightness: %v", err)
		}
		time.Sleep(5 * time.Millisecond) // Increase sleep time to reduce CPU usage
	}
	err = os.WriteFile(brightnessFile, []byte(fmt.Sprintf("%d", newBrightness)), 0644)
	if err != nil {
		log.Fatalf("Failed to set brightness: %v", err)
	}
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

// DetectBrightnessFiles automatically detects the brightness control files.
func detectBrightnessFiles() (string, string, error) {
	backlightDir := "/sys/class/backlight"
	files, err := filepath.Glob(filepath.Join(backlightDir, "*", "brightness"))
	if err != nil {
		return "", "", err
	}
	if len(files) == 0 {
		return "", "", fmt.Errorf("no brightness control files found")
	}
	brightnessFile := files[0]
	maxBrightnessFile := filepath.Join(filepath.Dir(brightnessFile), "max_brightness")
	return brightnessFile, maxBrightnessFile, nil
}

// DetectWebcamDevice automatically detects the webcam device.
func DetectWebcamDevice() (string, error) {
	devices, err := filepath.Glob("/dev/video*")
	if err != nil {
		return "", err
	}
	if len(devices) == 0 {
		return "", fmt.Errorf("no webcam devices found")
	}
	return devices[0], nil
}

// ParseDuration parses a duration string and converts it to a time.Duration.
func ParseDuration(duration string) (time.Duration, error) {
	duration = strings.ToLower(duration)
	if strings.HasSuffix(duration, "m") {
		duration = strings.TrimSuffix(duration, "m")
		duration += "m0s"
	} else if strings.HasSuffix(duration, "s") {
		duration = strings.TrimSuffix(duration, "s")
		duration += "s"
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 0, err
	}
	if d < 1*time.Second {
		return 0, fmt.Errorf("duration must be at least 5 seconds")
	}
	return d, nil
}

func main() {
	webcamDevice := flag.String("webcam", "", "Path to the webcam device")
	brightnessFile := flag.String("brightness", "", "Path to the brightness control file")
	maxBrightnessFile := flag.String("max-brightness", "", "Path to the max brightness control file")
	every := flag.String("every", "30s", "Time interval to capture a frame and adjust brightness")
	minBrightness := flag.Uint("min-brightness", 10, "Minimum brightness percentage (1-100)")
	singleSetBrightness := flag.Uint("set", 101, "Set brightness directly (1-100)")
	showMaxBrightness := flag.Bool("max", false, "Show maximum brightness value and exit")
	scaleFactor := flag.Int("scale", 120, "Scale factor for brightness transition")
	cmdInfo := &ccmd.CmdInfo{
		Name:        "glight",
		Authors:     []string{"xplshn"},
		Repository:  "https://github.com/xplshn/glight",
		Description: "Lets you controls your laptop's backlight easily",
		Synopsis:    "<|--webcam [filepath](/dev/video*)|--brightness [filepath](/sys/class/backlight/*/brightness)|--max-brightness [filepath](/sys/class/backlight/*/max_brightness)|--min-brightness [1-100](10)|--set [1-100]|--max [1-100]|--scale [1-100](120)> [FILE/s]",
	}
	helpPage, err := cmdInfo.GenerateHelpPage()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating help page:", err)
		os.Exit(1)
	}
	flag.Usage = func() {
		fmt.Print(helpPage)
	}
	flag.Parse()

	if *brightnessFile == "" || *maxBrightnessFile == "" {
		var err error
		*brightnessFile, *maxBrightnessFile, err = detectBrightnessFiles()
		if err != nil {
			log.Fatalf("Failed to detect brightness control files: %v", err)
		}
	}

	if *showMaxBrightness {
		maxBrightnessData, err := os.ReadFile(*maxBrightnessFile)
		if err != nil {
			log.Fatalf("Failed to read max brightness: %v", err)
		}
		fmt.Println(string(maxBrightnessData))
		os.Exit(0)
	}

	if *singleSetBrightness < 101 {
		setBrightness(float64(*singleSetBrightness), *brightnessFile, *maxBrightnessFile, uint8(*minBrightness), *scaleFactor)
		os.Exit(0)
	}

	if *webcamDevice == "" {
		var err error
		*webcamDevice, err = DetectWebcamDevice()
		if err != nil {
			log.Fatalf("Failed to detect webcam device: %v", err)
		}
	}

	interval, err := ParseDuration(*every)
	if err != nil {
		log.Fatalf("Invalid duration: %v", err)
	}

	for {
		cam, err := webcam.Open(*webcamDevice)
		if err != nil {
			log.Fatalf("Failed to open webcam: %v", err)
		}

		err = cam.StartStreaming()
		if err != nil {
			cam.Close()
			log.Fatalf("Failed to start streaming: %v", err)
		}

		err = cam.WaitForFrame(1) // Wait for at least one frame
		if err != nil {
			cam.StopStreaming()
			cam.Close()
			log.Printf("Failed to wait for frame: %v", err)
			continue
		}

		frame, err := cam.ReadFrame()
		if err != nil {
			cam.StopStreaming()
			cam.Close()
			log.Printf("Failed to read frame: %v", err)
			continue
		}

		img, _, err := image.Decode(bytes.NewReader(frame))
		if err != nil {
			cam.StopStreaming()
			cam.Close()
			log.Printf("Failed to decode frame: %v", err)
			continue
		}

		setBrightness(analyzeBrightness(img), *brightnessFile, *maxBrightnessFile, uint8(*minBrightness), *scaleFactor)

		cam.StopStreaming()
		cam.Close()

		time.Sleep(interval)
	}
}
