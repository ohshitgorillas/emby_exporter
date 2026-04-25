package emby

import (
	"strings"
)

type Status struct {
	ServerName string
	Hostname   string
	Version    string

	DeviceCount int

	MovieCount   int
	SeriesCount  int
	EpisodeCount int

	UserCount        int
	FailedLoginCount map[string]int

	MediaStreamCount   int
	TranscodingCount   int
	TranscodingBitRate map[string]int
	VideoBitRate       map[string]int
	AudioChannels      map[string]int
	AudioBitRate       map[string]int

	Sessions []SessionDetail
}

type SessionDetail struct {
	SessionId   string
	UserName    string
	Client      string
	DeviceName  string
	DeviceType  string
	AppVersion  string
	MediaType   string
	PlayMethod  string
	IsPaused    bool
	PositionSec float64
	RuntimeSec  float64
	ItemName    string
	ItemType    string

	Transcoding *TranscodeDetail
}

type TranscodeDetail struct {
	HwDecoder  string
	HwEncoder  string
	VideoCodec string
	AudioCodec string
	Container  string
	Reasons    string
	Bitrate    int
	Completion float64
	CpuUsage   float64
	Framerate  float64
	Width      int
	Height     int
}

func (c *Client) Status() (*Status, error) {
	s := &Status{
		FailedLoginCount:   make(map[string]int),
		TranscodingBitRate: make(map[string]int),
		VideoBitRate:       make(map[string]int),
		AudioChannels:      make(map[string]int),
		AudioBitRate:       make(map[string]int),
	}

	sys, err := c.GetSystemInfo()
	if err != nil {
		return nil, err
	}

	s.ServerName = strings.ToLower(sys.ServerName)
	s.Hostname = sys.WanAddress
	s.Version = sys.Version

	devices, err := c.GetDevices()
	if err != nil {
		return nil, err
	}

	s.DeviceCount = devices.TotalCount

	libary, err := c.GetLibraryItems()
	if err != nil {
		return nil, err
	}

	s.MovieCount = libary.MovieCount
	s.SeriesCount = libary.SeriesCount
	s.EpisodeCount = libary.EpisodeCount

	users, err := c.GetUsers()
	if err != nil {
		return nil, err
	}

	s.UserCount = users.TotalCount
	for _, user := range users.Users {
		s.FailedLoginCount[strings.ToLower(user.Name)] = user.Policy.InvalidLoginAttemptCount
	}

	sessions, err := c.GetSessions()
	if err != nil {
		return nil, err
	}

	for _, session := range *sessions {
		user := strings.ToLower(session.UserName)

		if session.TranscodingInfo != nil {
			s.TranscodingCount++
			s.TranscodingBitRate[user] += int(session.TranscodingInfo.Bitrate)
		}

		if len(session.NowPlayingItem.MediaStreams) > 0 {
			for _, stream := range session.NowPlayingItem.MediaStreams {
				if stream.Type == "Video" {
					s.MediaStreamCount++
					s.VideoBitRate[user] += int(stream.BitRate)
				}

				if stream.Type == "Audio" {
					s.AudioChannels[user] += int(stream.Channels)
					s.AudioBitRate[user] += int(stream.BitRate)
				}
			}
		}

		// Per-session detail (only for sessions with active playback).
		if session.NowPlayingItem.Id == "" {
			continue
		}

		detail := SessionDetail{
			SessionId:  session.Id,
			UserName:   user,
			Client:     session.Client,
			DeviceName: session.DeviceName,
			DeviceType: session.DeviceType,
			AppVersion: session.ApplicationVersion,
			MediaType:  session.NowPlayingItem.MediaType,
			ItemName:   session.NowPlayingItem.Name,
			ItemType:   session.NowPlayingItem.Type,
			RuntimeSec: float64(session.NowPlayingItem.RunTimeTicks) / 1e7,
		}
		if session.PlayState != nil {
			detail.PlayMethod = session.PlayState.PlayMethod
			detail.IsPaused = session.PlayState.IsPaused
			detail.PositionSec = float64(session.PlayState.PositionTicks) / 1e7
		}
		if t := session.TranscodingInfo; t != nil {
			detail.Transcoding = &TranscodeDetail{
				HwDecoder:  t.VideoDecoderHwAccel,
				HwEncoder:  t.VideoEncoderHwAccel,
				VideoCodec: t.VideoCodec,
				AudioCodec: t.AudioCodec,
				Container:  t.Container,
				Reasons:    strings.Join(t.TranscodeReasons, ","),
				Bitrate:    int(t.Bitrate),
				Completion: float64(t.CompletionPercentage),
				CpuUsage:   float64(t.CurrentCpuUsage),
				Framerate:  float64(t.Framerate),
				Width:      int(t.Width),
				Height:     int(t.Height),
			}
		}
		s.Sessions = append(s.Sessions, detail)
	}

	return s, nil
}
