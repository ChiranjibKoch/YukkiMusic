# 🎨 Thumbnail Customization Tutorial

This guide explains how to customize thumbnails that are displayed when playing songs in YukkiMusic bot.

## 📖 Overview

By default, YukkiMusic now processes thumbnails before sending them to your chat. The processed thumbnails include:
- **Original artwork** from YouTube/source
- **Text overlay** with song title and duration
- **Gradient background** for better text visibility
- **High-quality JPEG output** optimized for Telegram

## ⚙️ Configuration

### Environment Variables

Add these variables to your `.env` file to customize thumbnail behavior:

```env
# Enable/Disable thumbnail overlay (default: true)
THUMBNAIL_OVERLAY=true

# Custom font path for thumbnail text (optional)
THUMBNAIL_FONT=/path/to/your/custom/font.ttf
```

### Configuration Options

#### `THUMBNAIL_OVERLAY`
- **Type:** Boolean (`true` or `false`)
- **Default:** `true`
- **Description:** Controls whether to add text overlay on thumbnails
  - `true`: Thumbnails will have title and duration overlaid on the image
  - `false`: Original thumbnails will be used without any modification

**Example:**
```env
THUMBNAIL_OVERLAY=true
```

#### `THUMBNAIL_FONT`
- **Type:** String (file path)
- **Default:** System font (`/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf`)
- **Description:** Path to a custom TrueType Font (TTF) file for overlay text

**Example:**
```env
THUMBNAIL_FONT=/usr/share/fonts/truetype/custom/MyFont.ttf
```

## 🎯 Features

### 1. Automatic Text Overlay

When `THUMBNAIL_OVERLAY=true`, the bot automatically:
- Downloads the original thumbnail
- Adds a semi-transparent gradient at the bottom
- Overlays the song title (with automatic text wrapping)
- Displays the duration in the bottom-right corner
- Saves the processed image for Telegram

### 2. Smart Text Wrapping

Long song titles are automatically wrapped to fit within the thumbnail:
- Maximum of 2 lines
- Ellipsis (`...`) added if title is too long
- Dynamic font sizing based on image dimensions

### 3. Caching System

Processed thumbnails are cached for 30 minutes to improve performance:
- Same song requested multiple times uses cached thumbnail
- Reduces processing time and bandwidth
- Automatic cleanup of old cached files every 30 minutes

### 4. Fallback Handling

If thumbnail processing fails for any reason:
- Original thumbnail URL is used automatically
- No interruption to playback
- Error logged for debugging

## 📝 Examples

### Example 1: Enable Overlay with Default Font

```env
THUMBNAIL_OVERLAY=true
```

**Result:** Thumbnails will have white text overlay with song title and duration using the default system font.

### Example 2: Disable Overlay (Original Thumbnails)

```env
THUMBNAIL_OVERLAY=false
```

**Result:** Original thumbnails from YouTube/source will be sent without any modifications.

### Example 3: Custom Font

```env
THUMBNAIL_OVERLAY=true
THUMBNAIL_FONT=/usr/share/fonts/truetype/ubuntu/Ubuntu-Bold.ttf
```

**Result:** Thumbnails will use the Ubuntu Bold font for text overlay.

## 🎨 Customization Tips

### Installing Custom Fonts

To use a custom font on your VPS:

1. **Upload your TTF font file to the server:**
   ```bash
   mkdir -p /opt/yukki/fonts
   # Upload your font.ttf to this directory
   ```

2. **Set the font path in `.env`:**
   ```env
   THUMBNAIL_FONT=/opt/yukki/fonts/YourFont.ttf
   ```

3. **Restart the bot:**
   ```bash
   # Stop current instance
   # Restart: ./app
   ```

### Finding System Fonts

YukkiMusic automatically searches for fonts in common locations on different systems. If no custom font is specified, it will try these paths in order:

**Linux (Debian/Ubuntu):**
- `/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf`
- `/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf`
- `/usr/share/fonts/truetype/ubuntu/Ubuntu-Bold.ttf`

**Linux (Fedora/RHEL/CentOS):**
- `/usr/share/fonts/dejavu/DejaVuSans-Bold.ttf`
- `/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf`

**macOS:**
- `/System/Library/Fonts/Helvetica.ttc`

To list available system fonts on Linux:
```bash
fc-list | grep -i ttf
```

Common font locations:
- Ubuntu/Debian: `/usr/share/fonts/truetype/`
- CentOS/RHEL: `/usr/share/fonts/`
- Custom fonts: `/usr/local/share/fonts/`

## 🖼️ Technical Details

### Image Processing

- **Library Used:** `github.com/fogleman/gg` for drawing, `github.com/nfnt/resize` for resizing
- **Output Format:** JPEG with 85% quality
- **Maximum Dimensions:** 1280x720 (maintains aspect ratio)
- **Overlay Gradient:** Semi-transparent black (0, 0, 0, 180) covering bottom 25% of image

### Text Rendering

- **Title Font Size:** Dynamic, calculated as `image_width / 25`
- **Duration Font Size:** Dynamic, calculated as `image_width / 35`
- **Text Color:** White (255, 255, 255, 255)
- **Shadow:** Black (0, 0, 0, 200) offset by 1-2 pixels
- **Title Position:** Center-bottom
- **Duration Position:** Bottom-right corner

### Storage & Cleanup

- **Temporary Location:** System temp directory (platform-specific)
  - Linux/macOS: `/tmp/yukki_thumbnails/`
  - Windows: `%TEMP%\yukki_thumbnails\`
- **Filename Format:** `thumb_<timestamp>.jpg`
- **Cleanup Interval:** Every 30 minutes
- **Cache Retention:** Files older than 1 hour are automatically deleted

## 🔧 Troubleshooting

### Issue: Thumbnails not showing overlay

**Solution:**
1. Check that `THUMBNAIL_OVERLAY=true` in your `.env`
2. Ensure the bot has write permissions to system temp directory
3. Check logs for font loading errors
4. Verify that at least one system font is available on your system

### Issue: Custom font not working

**Solution:**
1. Verify the font file exists at the specified path
2. Ensure the file is a valid TTF font
3. Check file permissions (should be readable by the bot user)
4. Try using an absolute path

### Issue: "Failed to process thumbnail" errors

**Solution:**
1. Ensure you have internet connectivity to download thumbnails
2. Check if the source thumbnail URL is accessible
3. Verify sufficient disk space in `/tmp`
4. The bot will automatically fall back to original thumbnails

## 📊 Performance Considerations

- **First Request:** ~2-3 seconds to download and process thumbnail
- **Cached Requests:** Instant (no processing needed)
- **Disk Usage:** Minimal (~1-5 MB depending on cache size)
- **Memory Usage:** Negligible (images processed on-demand)

## 🎓 Advanced Customization

For developers who want to customize the thumbnail processing further, edit:
- **File:** `internal/utils/thumbnail.go`
- **Key Functions:**
  - `ProcessThumbnail()` - Main processing function
  - `addOverlay()` - Overlay rendering logic
  - `DefaultThumbnailConfig()` - Default settings

You can modify:
- Background colors
- Text colors and styles
- Gradient opacity and size
- Font sizes
- Image quality and dimensions
- Text positioning

## 📖 Related Documentation

- [README.md](.github/README.md) - Main documentation
- [sample.env](sample.env) - Configuration template
- Platform System: [internal/platforms/README.md](internal/platforms/README.md)

## 💡 Tips

1. **Best Font Choice:** Use bold fonts for better readability on images
2. **Performance:** Keep `THUMBNAIL_OVERLAY=true` for better user experience
3. **Quality:** Default settings are optimized for Telegram - don't change unless needed
4. **Caching:** The 30-minute cache greatly improves performance for popular songs

## 🤝 Support

If you need help with thumbnail customization:
- **Telegram:** [@TheTeamVk](https://t.me/TheTeamVk)
- **Issues:** [GitHub Issues](https://github.com/TheTeamVivek/YukkiMusic/issues)

---

**Note:** Thumbnail customization is available from YukkiMusic v2.0 onwards.
