/*
 * This file is part of YukkiMusic.
 *
 * YukkiMusic — A Telegram bot that streams music into group voice chats with seamless playback and control.
 * Copyright (C) 2025 TheTeamVivek
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */
package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Laky-64/gologging"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"

	"github.com/TheTeamVivek/YukkiMusic/config"
)

var thumbnailCache = NewCache[string, string](30 * time.Minute)
var logger = gologging.GetLogger("Thumbnail")

// Common font paths for different systems
var defaultFontPaths = []string{
	"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",        // Debian/Ubuntu
	"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",             // Debian/Ubuntu
	"/usr/share/fonts/dejavu/DejaVuSans-Bold.ttf",                 // Fedora/RHEL
	"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf", // Alternative
	"/System/Library/Fonts/Helvetica.ttc",                         // macOS
	"/usr/share/fonts/truetype/ubuntu/Ubuntu-Bold.ttf",            // Ubuntu
}

// ThumbnailConfig holds configuration for thumbnail customization
type ThumbnailConfig struct {
	AddOverlay      bool
	OverlayText     string
	TitleText       string
	DurationText    string
	BackgroundColor color.Color
	TextColor       color.Color
	Width           uint
	Height          uint
	Quality         int
}

// DefaultThumbnailConfig returns the default configuration
func DefaultThumbnailConfig() *ThumbnailConfig {
	return &ThumbnailConfig{
		AddOverlay:      config.ThumbnailOverlay,
		BackgroundColor: color.RGBA{0, 0, 0, 180}, // Semi-transparent black
		TextColor:       color.RGBA{255, 255, 255, 255}, // White
		Width:           1280,
		Height:          720,
		Quality:         85,
	}
}

// ProcessThumbnail downloads and processes a thumbnail with custom overlay
func ProcessThumbnail(thumbnailURL, title, duration string) (string, error) {
	if thumbnailURL == "" {
		return "", nil
	}

	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s:%s", thumbnailURL, title, duration)
	if cached, ok := thumbnailCache.Get(cacheKey); ok {
		if _, err := os.Stat(cached); err == nil {
			return cached, nil
		}
	}

	// Download original thumbnail
	img, err := downloadImage(thumbnailURL)
	if err != nil {
		logger.ErrorF("Failed to download thumbnail: %v", err)
		return "", err
	}

	// If overlay is disabled, just save and return original
	if !config.ThumbnailOverlay {
		outputPath := generateOutputPath()
		if err := saveImage(img, outputPath, 95); err != nil {
			return "", err
		}
		thumbnailCache.Set(cacheKey, outputPath)
		return outputPath, nil
	}

	// Create custom thumbnail with overlay
	cfg := DefaultThumbnailConfig()
	cfg.TitleText = title
	cfg.DurationText = duration

	processedImg, err := addOverlay(img, cfg)
	if err != nil {
		logger.ErrorF("Failed to add overlay: %v", err)
		return "", err
	}

	// Save processed image
	outputPath := generateOutputPath()
	if err := saveImage(processedImg, outputPath, cfg.Quality); err != nil {
		return "", err
	}

	thumbnailCache.Set(cacheKey, outputPath)
	return outputPath, nil
}

// downloadImage downloads an image from a URL
func downloadImage(url string) (image.Image, error) {
	// Clean URL
	url = CleanURL(url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}

// loadFont attempts to load a font, trying custom path first, then default paths
func loadFont(dc *gg.Context, fontSize float64) error {
	// Try custom font first
	customFont := config.ThumbnailFont
	if customFont != "" && fileExists(customFont) {
		if err := dc.LoadFontFace(customFont, fontSize); err == nil {
			return nil
		}
		logger.WarnF("Failed to load custom font %s: trying defaults", customFont)
	}

	// Try default system fonts
	for _, fontPath := range defaultFontPaths {
		if fileExists(fontPath) {
			if err := dc.LoadFontFace(fontPath, fontSize); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("no suitable font found")
}

// addOverlay adds text overlay to the thumbnail
func addOverlay(img image.Image, cfg *ThumbnailConfig) (image.Image, error) {
	// Resize image if needed
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())

	if width > cfg.Width || height > cfg.Height {
		img = resize.Resize(cfg.Width, cfg.Height, img, resize.Lanczos3)
		bounds = img.Bounds()
		width = uint(bounds.Dx())
		height = uint(bounds.Dy())
	}

	// Create a new context for drawing
	dc := gg.NewContextForImage(img)

	// Add gradient overlay at bottom for better text visibility
	gradientHeight := float64(height) * 0.25
	for i := 0.0; i < gradientHeight; i++ {
		alpha := uint8((i / gradientHeight) * 180)
		dc.SetColor(color.RGBA{0, 0, 0, alpha})
		dc.DrawRectangle(0, float64(height)-gradientHeight+i, float64(width), 1)
		dc.Fill()
	}

	// Calculate font size based on image dimensions
	titleFontSize := float64(width) / 25.0
	durationFontSize := float64(width) / 35.0

	// Load font for title
	if err := loadFont(dc, titleFontSize); err != nil {
		logger.WarnF("Failed to load font, returning original image: %v", err)
		return img, nil
	}

	// Draw title text
	if cfg.TitleText != "" {
		// Wrap text if too long
		maxWidth := float64(width) * 0.9
		wrappedTitle := wrapText(dc, cfg.TitleText, maxWidth)

		// Draw text with shadow for better visibility
		dc.SetColor(color.RGBA{0, 0, 0, 200}) // Shadow
		x := float64(width) / 2
		y := float64(height) - 60
		dc.DrawStringAnchored(wrappedTitle, x+2, y+2, 0.5, 0.5)

		dc.SetColor(cfg.TextColor) // Main text
		dc.DrawStringAnchored(wrappedTitle, x, y, 0.5, 0.5)
	}

	// Draw duration text
	if cfg.DurationText != "" {
		// Load font for duration
		if err := loadFont(dc, durationFontSize); err != nil {
			logger.WarnF("Failed to load font for duration: %v", err)
		} else {
			dc.SetColor(color.RGBA{0, 0, 0, 200}) // Shadow
			x := float64(width) - 20
			y := float64(height) - 20
			dc.DrawStringAnchored(cfg.DurationText, x+1, y+1, 1.0, 1.0)

			dc.SetColor(color.RGBA{255, 255, 255, 255}) // Main text
			dc.DrawStringAnchored(cfg.DurationText, x, y, 1.0, 1.0)
		}
	}

	return dc.Image(), nil
}

// wrapText wraps text to fit within maxWidth
func wrapText(dc *gg.Context, text string, maxWidth float64) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		w, _ := dc.MeasureString(testLine)
		if w > maxWidth && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// Limit to 2 lines and add ellipsis if needed
	if len(lines) > 2 {
		lines = lines[:2]
		lastLine := lines[1]
		// Truncate last line and add ellipsis
		for {
			w, _ := dc.MeasureString(lastLine + "...")
			if w <= maxWidth {
				lines[1] = lastLine + "..."
				break
			}
			words := strings.Fields(lastLine)
			if len(words) <= 1 {
				break
			}
			lastLine = strings.Join(words[:len(words)-1], " ")
		}
	}

	return strings.Join(lines, "\n")
}

// saveImage saves an image to disk
func saveImage(img image.Image, path string, quality int) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Save as JPEG with specified quality
	if strings.HasSuffix(strings.ToLower(path), ".png") {
		return png.Encode(file, img)
	}

	return jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
}

// generateOutputPath generates a unique output path for processed thumbnails
func generateOutputPath() string {
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("thumb_%d.jpg", timestamp)
	return filepath.Join(os.TempDir(), "yukki_thumbnails", filename)
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CleanupOldThumbnails removes old thumbnail files
func CleanupOldThumbnails() {
	dir := filepath.Join(os.TempDir(), "yukki_thumbnails")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		logger.ErrorF("Failed to read thumbnail directory: %v", err)
		return
	}

	now := time.Now()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// Remove files older than 1 hour
		if now.Sub(info.ModTime()) > time.Hour {
			path := filepath.Join(dir, file.Name())
			if err := os.Remove(path); err != nil {
				logger.ErrorF("Failed to remove old thumbnail %s: %v", path, err)
			}
		}
	}
}

// FormatDurationForThumbnail formats duration for thumbnail overlay
func FormatDurationForThumbnail(seconds int) string {
	if seconds <= 0 {
		return "00:00"
	}

	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}
