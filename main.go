package main

import (
	"fmt"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Brightness Control App")

	displayDropdown := widget.NewSelect([]string{}, func(displayName string) {
		fmt.Printf("Selected display: %s\n", displayName)
	})

	populateDisplays(displayDropdown)

	brightnessSlider := widget.NewSlider(10, 100)
	brightnessLabel := widget.NewLabel("Brightness: 0%")

	brightnessSlider.OnChanged = func(value float64) {
		brightnessLabel.SetText(fmt.Sprintf("Brightness: %.0f%%", value))
		adjustBrightness(displayDropdown.Selected, value)
	}

	myWindow.SetContent(container.NewVBox(
		widget.NewLabel("Set brigness below hardware minimum"),
		layout.NewSpacer(),
		widget.NewLabel("Select Display:"),
		displayDropdown,
		widget.NewLabel("Adjust Brightness"),
		brightnessSlider,
		brightnessLabel,
		layout.NewSpacer(),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Quit", theme.CancelIcon(), func() {
				myApp.Quit()
			}),
		),
	))

	myWindow.ShowAndRun()
}

func populateDisplays(dropdown *widget.Select) {
	cmd := exec.Command("xrandr")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error fetching displays: %v\n", err)
		return
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, " connected") {
			parts := strings.Fields(line)
			displayName := parts[0]
			dropdown.Options = append(dropdown.Options, displayName)
		}
	}
}

func adjustBrightness(displayName string, value float64) {
	cmd := exec.Command("xrandr", "--output", displayName, "--brightness", fmt.Sprintf("%.2f", value/100))
	fmt.Printf("Setting brightness to %.2f\n", value/100)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error adjusting brightness: %v\n", err)
	}
}
