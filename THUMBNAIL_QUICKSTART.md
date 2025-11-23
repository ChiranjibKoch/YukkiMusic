# Thumbnail Customization - Quick Start Guide

## What's New?

YukkiMusic now automatically customizes song thumbnails before sending them! When you play a song, the bot will:

1. Download the original thumbnail from YouTube/source
2. Add a beautiful overlay with:
   - **Song Title** (center, with smart text wrapping)
   - **Duration** (bottom-right corner)
   - **Semi-transparent gradient** for better text visibility
3. Cache the processed thumbnail for future use
4. Send the customized thumbnail to your chat

## Before & After Example

**Before (Default YouTube Thumbnail):**
- Plain YouTube thumbnail with no additional information

**After (With YukkiMusic Customization):**
- Same thumbnail + elegant text overlay
- Song title clearly visible
- Duration shown in corner
- Professional gradient background

## Quick Setup

### Option 1: Use Default Settings (Recommended)

Just update your bot - thumbnail customization is **enabled by default**!

No configuration needed. The bot will automatically:
- Process all thumbnails
- Use system fonts (DejaVu, Liberation, etc.)
- Store files in system temp directory
- Clean up old files automatically

### Option 2: Disable Thumbnails (If You Prefer)

Add this to your `.env` file:

```env
THUMBNAIL_OVERLAY=false
```

This will use original thumbnails without any modification.

### Option 3: Use Custom Font

Add this to your `.env` file:

```env
THUMBNAIL_OVERLAY=true
THUMBNAIL_FONT=/path/to/your/font.ttf
```

Example with Ubuntu font:
```env
THUMBNAIL_FONT=/usr/share/fonts/truetype/ubuntu/Ubuntu-Bold.ttf
```

## How It Works

1. **User plays a song** using `/play` command
2. **Bot fetches** song info and original thumbnail
3. **Processing** (if overlay enabled):
   - Downloads thumbnail
   - Creates gradient overlay
   - Adds song title (wrapped if long)
   - Adds duration
   - Saves processed image
4. **Bot sends** the customized thumbnail with "Now Playing" message
5. **Caching**: Same song next time? Uses cached version instantly!

## Features

✅ **Automatic Processing** - No manual intervention needed  
✅ **Smart Caching** - Faster playback for popular songs  
✅ **Cross-Platform** - Works on Linux, macOS, Windows  
✅ **Fallback Safe** - Uses original if processing fails  
✅ **Configurable** - Enable/disable as you prefer  
✅ **Custom Fonts** - Use your favorite fonts  
✅ **Auto-Cleanup** - Old files removed automatically  

## Performance

- **First Play**: ~2-3 seconds (download + process)
- **Cached Play**: Instant (no processing)
- **Disk Usage**: Minimal (~1-5 MB)
- **Memory**: Negligible (processed on-demand)

## Common Questions

**Q: Will this slow down the bot?**  
A: First play of a song takes 2-3 seconds extra. Subsequent plays are instant (cached).

**Q: Can I disable it?**  
A: Yes! Set `THUMBNAIL_OVERLAY=false` in your `.env` file.

**Q: What if processing fails?**  
A: The bot automatically falls back to the original thumbnail. Music playback is never interrupted.

**Q: Do I need to install fonts?**  
A: No! The bot uses system fonts automatically. Custom fonts are optional.

**Q: How much storage does it use?**  
A: Very little (~1-5 MB). Old thumbnails are automatically cleaned up every 30 minutes.

## Troubleshooting

### Thumbnails not showing overlay

1. Check `.env` has `THUMBNAIL_OVERLAY=true`
2. Check bot logs for errors
3. Verify system has fonts installed:
   ```bash
   fc-list | grep -i ttf
   ```

### "Failed to load font" error

The bot tries multiple font paths automatically. If you see this error:
1. Install DejaVu fonts: `sudo apt install fonts-dejavu`
2. Or specify a custom font in `.env`

### Out of disk space

The bot cleans up old files automatically. But if needed:
```bash
# Remove thumbnail cache manually
rm -rf /tmp/yukki_thumbnails/
```

## Advanced Customization

For developers who want to customize further, edit `internal/utils/thumbnail.go`:

- Change gradient colors/opacity
- Modify text positioning
- Adjust font sizes
- Change image quality
- Customize text colors

See [THUMBNAIL_CUSTOMIZATION.md](THUMBNAIL_CUSTOMIZATION.md) for full documentation.

## Support

Need help? Reach out:
- **Telegram**: [@TheTeamVk](https://t.me/TheTeamVk)
- **Issues**: [GitHub Issues](https://github.com/TheTeamVivek/YukkiMusic/issues)

---

**Enjoy beautiful thumbnails with YukkiMusic! 🎨🎵**
