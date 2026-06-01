/*
 * ● YukkiMusic
 * ○ A high-performance engine for streaming music in Telegram voicechats.
 *
 * Copyright (C) 2026 TheTeamVivek
 *
 * This program is free software: you can redistribute it and/or modify it under the
 * terms of the GNU General Public License as published by the Free Software Foundation,
 * either version 3 of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
 * PARTICULAR PURPOSE. See the GNU General Public License for more details.
 *
 * Repository: https://github.com/TheTeamVivek/YukkiMusic
 */

package platforms

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Laky-64/gologging"
	"github.com/amarnathcjd/gogram/telegram"

	"main/internal/config"
	state "main/internal/core/models"
	"main/internal/utils"
)

const (
	PlatformArchTube      state.PlatformName = "ArchTube"
	archTubeAPIKeyHeader                     = "ArchFreeYt"
)

type archTubeMetaResponse struct {
	Title     string `json:"title"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	Author    string `json:"author"`
}

type archTubeSearchItem struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Duration int    `json:"duration"`
	Author   string `json:"author"`
	Thumb    string `json:"thumbnail"`
}

type ArchTubePlatform struct {
	name state.PlatformName
}

func init() {
	Register(75, &ArchTubePlatform{
		name: PlatformArchTube,
	})
}

func (a *ArchTubePlatform) Name() state.PlatformName {
	return a.name
}

func (a *ArchTubePlatform) CanGetTracks(query string) bool {
	if config.ArchTubeAPIURL == "" || config.ArchTubeAPIKey == "" {
		return false
	}
	return !youtubeLinkRegex.MatchString(query)
}

func (a *ArchTubePlatform) GetTracks(
	query string,
	_ bool,
) ([]*state.Track, error) {
	reqURL := fmt.Sprintf(
		"%s/api/search?q=%s&max=5",
		config.ArchTubeAPIURL,
		url.QueryEscape(query),
	)

	var results []archTubeSearchItem
	resp, err := rc.R().
		SetHeader(archTubeAPIKeyHeader, config.ArchTubeAPIKey).
		SetResult(&results).
		Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("archtube search request failed: %w",
			sanitizeAPIError(err, config.ArchTubeAPIKey))
	}
	if resp.IsError() {
		return nil, sanitizeAPIError(fmt.Errorf(
			"archtube search failed with status %d: %s",
			resp.StatusCode(), resp.String(),
		), config.ArchTubeAPIKey)
	}
	if len(results) == 0 {
		return nil, errors.New("archtube: no results found")
	}

	tracks := make([]*state.Track, 0, len(results))
	for _, item := range results {
		tracks = append(tracks, &state.Track{
			Title:    item.Title,
			URL:      item.URL,
			Duration: item.Duration,
			Artwork:  item.Thumb,
			Source:   PlatformArchTube,
		})
	}
	return tracks, nil
}

func (a *ArchTubePlatform) CanDownload(source state.PlatformName) bool {
	if config.ArchTubeAPIURL == "" || config.ArchTubeAPIKey == "" {
		return false
	}
	return source == PlatformYouTube || source == PlatformArchTube
}

func (a *ArchTubePlatform) Download(
	ctx context.Context,
	track *state.Track,
	statusMsg *telegram.NewMessage,
) (string, error) {
	if cached := findFile(track); cached != "" {
		gologging.Debug("ArchTube: Download -> Cached File -> " + cached)
		return cached, nil
	}

	var pm *telegram.ProgressManager
	if statusMsg != nil {
		pm = utils.GetProgress(statusMsg)
	}
	_ = pm

	var (
		endpoint string
		ext      string
	)
	if track.Video {
		endpoint = "download/video"
		ext = ".mp4"
	} else {
		endpoint = "download/audio"
		ext = ".mp3"
	}

	dlURL := fmt.Sprintf(
		"%s/%s?url=%s",
		config.ArchTubeAPIURL,
		endpoint,
		url.QueryEscape(track.URL),
	)

	path := getPath(track, ext)

	resp, err := rc.R().
		SetContext(ctx).
		SetHeader(archTubeAPIKeyHeader, config.ArchTubeAPIKey).
		SetOutputFileName(path).
		Get(dlURL)
	if err != nil {
		os.Remove(path)
		if errors.Is(err, context.Canceled) ||
			errors.Is(err, context.DeadlineExceeded) {
			return "", err
		}
		return "", fmt.Errorf("archtube download failed: %w",
			sanitizeAPIError(err, config.ArchTubeAPIKey))
	}

	if resp.IsError() {
		os.Remove(path)
		return "", sanitizeAPIError(fmt.Errorf(
			"archtube download failed with status %d: %s",
			resp.StatusCode(), resp.String(),
		), config.ArchTubeAPIKey)
	}

	if !fileExists(path) {
		return "", errors.New("archtube: empty file returned by API")
	}

	gologging.Debug("ArchTube: Download -> " + path)
	return path, nil
}

func (a *ArchTubePlatform) fetchMeta(
	ctx context.Context,
	mediaURL string,
) (*archTubeMetaResponse, error) {
	type metaRequest struct {
		URL string `json:"url"`
	}

	var meta archTubeMetaResponse
	resp, err := rc.R().
		SetContext(ctx).
		SetHeader(archTubeAPIKeyHeader, config.ArchTubeAPIKey).
		SetBody(metaRequest{URL: mediaURL}).
		SetResult(&meta).
		Post(config.ArchTubeAPIURL + "/meta")
	if err != nil {
		return nil, fmt.Errorf("archtube meta request failed: %w",
			sanitizeAPIError(err, config.ArchTubeAPIKey))
	}
	if resp.IsError() {
		return nil, sanitizeAPIError(fmt.Errorf(
			"archtube meta failed with status %d: %s",
			resp.StatusCode(), resp.String(),
		), config.ArchTubeAPIKey)
	}
	return &meta, nil
}
