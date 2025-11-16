package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Laky-64/gologging"
	xdraw "golang.org/x/image/draw"
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

	// Fetch artwork
	resp, err := http.Get(CleanURL(track.Artwork))
	if err != nil {
		gologging.ErrorF("Failed to fetch artwork: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		gologging.ErrorF("Artwork returned status %d", resp.StatusCode)
		return ""
	}

	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		gologging.ErrorF("Failed reading artwork body: %v", err)
		return ""
	}

	contentType := http.DetectContentType(imgData)

	var src image.Image

	switch contentType {
	case "image/jpeg":
		src, err = jpeg.Decode(bytes.NewReader(imgData))
	case "image/png":
		src, err = png.Decode(bytes.NewReader(imgData))
	default:
		gologging.ErrorF("Unsupported artwork format: %s", contentType)
		return ""
	}

	if err != nil {
		gologging.ErrorF("Failed to decode artwork: %v", err)
		return ""
	}

	// Canvas
	const W, H = 1920, 1080
	base := image.NewRGBA(image.Rect(0, 0, W, H))

	// Background color
	bgColor := color.RGBA{18, 27, 33, 255}
	draw.Draw(base, base.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Wave area / bottom bar
	waveColor := color.RGBA{28, 37, 45, 255}
	waveRect := image.Rect(0, H-400, W, H)
	draw.Draw(base, waveRect, &image.Uniform{waveColor}, image.Point{}, draw.Over)

	// Album resized
	album := image.NewRGBA(image.Rect(0, 0, 650, 650))
	xdraw.CatmullRom.Scale(album, album.Bounds(), src, src.Bounds(), xdraw.Over, nil)

	draw.Draw(base, image.Rect(180, 220, 830, 870), album, image.Point{}, draw.Over)

	// Text drawing
	face := basicfont.Face7x13
	drawer := &font.Drawer{
		Dst:  base,
		Face: face,
	}

	// Playing
	drawer.Src = image.NewUniform(color.RGBA{185, 192, 199, 255})
	drawer.Dot = fixed.Point26_6{X: fixed.I(900), Y: fixed.I(330)}
	drawer.DrawString("Playing")

	// Title
	drawer.Src = image.White
	drawer.Dot = fixed.Point26_6{X: fixed.I(900), Y: fixed.I(420)}
	drawer.DrawString(title)

	// Artist
	drawer.Src = image.NewUniform(color.RGBA{205, 205, 205, 255})
	drawer.Dot = fixed.Point26_6{X: fixed.I(900), Y: fixed.I(550)}
	drawer.DrawString(artist)

	// Duration
	drawer.Src = image.NewUniform(color.RGBA{180, 180, 180, 255})
	drawer.Dot = fixed.Point26_6{X: fixed.I(900), Y: fixed.I(650)}
	drawer.DrawString(fmt.Sprintf("Duration: %d", duration))

	// Save PNG
	outFile, err := os.Create(cachePath)
	if err != nil {
		return ""
	}
	defer outFile.Close()

	png.Encode(outFile, base)

	return cachePath
}
