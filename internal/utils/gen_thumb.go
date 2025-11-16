package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Laky-64/gologging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"main/internal/state"
)

const cacheDir = "cache"

func GenThumb(track *state.Track) string {
	if track == nil || track.Artwork == "" {
		return ""
	}

	os.MkdirAll(cacheDir, 0o755)

	cachePath := filepath.Join(cacheDir, fmt.Sprintf("%s.png", track.ID))

	if _, err := os.Stat(cachePath); err == nil {
		return cachePath
	}

	title := track.Title
	artist := "Vivek"
	duration := track.Duration

	resp, err := http.Get(track.Artwork)
	if err != nil {
		gologging.ErrorF("Failed to get artwork %v", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		gologging.ErrorF("Failed to get artwork StatusCode %d", resp.StatusCode)

		return ""
	}
	imgData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		gologging.ErrorF("Failed to read raw image %v", err)

		return ""
	}

	thumbPath := filepath.Join(cacheDir, fmt.Sprintf("thumb_%s.jpg", track.ID))
	err = ioutil.WriteFile(thumbPath, imgData, 0o644)
	if err != nil {
		gologging.ErrorF("Failed to write raw artwork %v", err)

		return ""
	}

	file, err := os.Open(thumbPath)
	if err != nil {
		gologging.ErrorF("Failed to open thumbPath %v", err)

		return ""
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		gologging.ErrorF("Failed to decode jpeg %v", err)

		return ""
	}

	const W, H = 1920, 1080
	base := image.NewRGBA(image.Rect(0, 0, W, H))
	bgColor := color.RGBA{18, 27, 33, 255}
	draw.Draw(base, base.Bounds(), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

	waveColor := color.RGBA{28, 37, 45, 255}
	waveRect := image.Rect(0, H-400, W, H)
	draw.Draw(base, waveRect, &image.Uniform{C: waveColor}, image.Point{}, draw.Over)

	album := image.NewRGBA(image.Rect(0, 0, 650, 650))
	draw.CatmullRom.Scale(album, album.Bounds(), img, img.Bounds(), draw.Over, nil)

	draw.Draw(base, image.Rect(180, 220, 180+650, 220+650), album, image.Point{}, draw.Over)

	face := basicfont.Face7x13

	drawer := &font.Drawer{
		Dst:  base,
		Src:  image.NewUniform(color.RGBA{185, 192, 199, 255}), // light grey color
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(900),
			Y: fixed.I(330),
		},
	}

	// Draw "Playing"
	drawer.DrawString("Playing")

	// Draw the track title (white color)
	drawer.Src = image.NewUniform(color.White)
	drawer.Dot = fixed.Point26_6{X: fixed.I(900), Y: fixed.I(420)}
	drawer.DrawString(title)

	// Draw the artist name (light grey)
	drawer.Src = image.NewUniform(color.RGBA{205, 205, 205, 255})
	drawer.Dot = fixed.Point26_6{X: fixed.I(900), Y: fixed.I(550)}
	drawer.DrawString(artist)

	// Draw the duration (grey)
	drawer.Src = image.NewUniform(color.RGBA{180, 180, 180, 255})
	drawer.Dot = fixed.Point26_6{X: fixed.I(900), Y: fixed.I(650)}
	drawer.DrawString(fmt.Sprintf("Duration: %d", duration))

	os.Remove(thumbPath)
	outFile, err := os.Create(cachePath)
	if err != nil {
		return ""
	}
	defer outFile.Close()
	png.Encode(outFile, base)

	return cachePath
}
