package collector

import (
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/ohshitgorillas/emby_exporter/emby"
)

var _ StatusSource = &emby.Client{}

type StatusSource interface {
	Status() (*emby.Status, error)
}

type Collector struct {
	Info *prometheus.Desc

	DeviceCount        *prometheus.Desc
	MovieCount         *prometheus.Desc
	SeriesCount        *prometheus.Desc
	EpisodeCount       *prometheus.Desc
	UserCount          *prometheus.Desc
	FailedLoginCount   *prometheus.Desc
	MediaStreamCount   *prometheus.Desc
	TranscodingCount   *prometheus.Desc
	TranscodingBitRate *prometheus.Desc
	VideoBitRate       *prometheus.Desc
	AudioBitRate       *prometheus.Desc
	AudioChannels      *prometheus.Desc

	LibraryItemCount    *prometheus.Desc

	Up             *prometheus.Desc
	ScrapeDuration *prometheus.Desc

	ScheduledTaskInfo            *prometheus.Desc
	ScheduledTaskLastDuration    *prometheus.Desc
	ScheduledTaskLastRunTimestamp *prometheus.Desc
	PluginUpdateAvailable        *prometheus.Desc

	SessionInfo         *prometheus.Desc
	SessionPositionSec  *prometheus.Desc
	SessionRuntimeSec   *prometheus.Desc
	TranscodeInfo       *prometheus.Desc
	TranscodeBitRate    *prometheus.Desc
	TranscodeCompletion *prometheus.Desc
	TranscodeCpuUsage   *prometheus.Desc
	TranscodeFramerate  *prometheus.Desc
	TranscodeWidth      *prometheus.Desc
	TranscodeHeight     *prometheus.Desc

	ss StatusSource
}

func New(ss StatusSource) prometheus.Collector {
	labels := []string{"server"}
	namespace := "emby"

	return &Collector{
		Info: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "info"),
			"Metadata about a given Emby server.",
			[]string{"server", "hostname", "version"},
			nil,
		),

		DeviceCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "device", "count"),
			"Number of devices configured in Emby.",
			labels,
			nil,
		),

		MovieCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "movie", "count"),
			"Number of movies available in Emby.",
			labels,
			nil,
		),

		SeriesCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "series", "count"),
			"Number of tv shows available in Emby.",
			labels,
			nil,
		),

		EpisodeCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "episode", "count"),
			"Number of tv show episodes available in Emby.",
			labels,
			nil,
		),

		UserCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "user", "count"),
			"Number of users configured in Emby.",
			labels,
			nil,
		),

		FailedLoginCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "failed", "login"),
			"Failed login counts per user in Emby.",
			[]string{"server", "user"},
			nil,
		),

		MediaStreamCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "stream", "count"),
			"Number of media streams being handled by Emby.",
			labels,
			nil,
		),

		TranscodingCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcoding", "count"),
			"Number of media streams being transcoded by Emby.",
			labels,
			nil,
		),

		TranscodingBitRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcoding", "bitrate"),
			"Bitrate of transcoded mediastream.",
			[]string{"server", "user"},
			nil,
		),

		VideoBitRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "video", "bitrate"),
			"Bitrate of original video file.",
			[]string{"server", "user"},
			nil,
		),

		AudioBitRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "audio", "bitrate"),
			"Bitrate of original audio file.",
			[]string{"server", "user"},
			nil,
		),

		AudioChannels: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "audio", "channels"),
			"Number of channels in original audio file.",
			[]string{"server", "user"},
			nil,
		),

		Up: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "", "up"),
			"1 if the last scrape of the Emby API succeeded, 0 otherwise.",
			nil,
			nil,
		),

		ScrapeDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scrape", "duration_seconds"),
			"Duration of the last scrape of the Emby API in seconds.",
			nil,
			nil,
		),

		LibraryItemCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "library", "item_count"),
			"Number of items in a library, broken down by item type.",
			[]string{"server", "library", "collection_type", "type"},
			nil,
		),

		ScheduledTaskInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "info"),
			"Scheduled task metadata. Always 1.",
			[]string{"server", "task", "category", "state", "last_status"},
			nil,
		),

		ScheduledTaskLastDuration: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "last_duration_seconds"),
			"Duration in seconds of the most recent run of a scheduled task.",
			[]string{"server", "task"},
			nil,
		),

		ScheduledTaskLastRunTimestamp: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "scheduled_task", "last_run_timestamp_seconds"),
			"Unix timestamp of the most recent run of a scheduled task.",
			[]string{"server", "task"},
			nil,
		),

		PluginUpdateAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "plugin", "update_available"),
			"Available update for an installed plugin. Always 1.",
			[]string{"server", "name", "version"},
			nil,
		),

		SessionInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "info"),
			"Active playback session metadata. Always 1.",
			[]string{"server", "session_id", "user", "client", "device", "device_type", "app_version", "media_type", "play_method", "is_paused", "item_type"},
			nil,
		),

		SessionPositionSec: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "position_seconds"),
			"Current playback position in seconds.",
			[]string{"server", "session_id", "user"},
			nil,
		),

		SessionRuntimeSec: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "session", "runtime_seconds"),
			"Total runtime of the now-playing item in seconds.",
			[]string{"server", "session_id", "user"},
			nil,
		),

		TranscodeInfo: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcode", "info"),
			"Active transcode session metadata. Always 1.",
			[]string{"server", "session_id", "user", "hw_decoder", "hw_encoder", "video_codec", "audio_codec", "container", "reasons"},
			nil,
		),

		TranscodeBitRate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcode", "bitrate"),
			"Bitrate of an active transcode in bits per second.",
			[]string{"server", "session_id", "user"},
			nil,
		),

		TranscodeCompletion: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcode", "completion_percent"),
			"Completion percentage of an active transcode.",
			[]string{"server", "session_id", "user"},
			nil,
		),

		TranscodeCpuUsage: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcode", "cpu_usage"),
			"Current CPU usage of an active transcode (Emby-reported, fractional).",
			[]string{"server", "session_id", "user"},
			nil,
		),

		TranscodeFramerate: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcode", "framerate"),
			"Current framerate of an active transcode.",
			[]string{"server", "session_id", "user"},
			nil,
		),

		TranscodeWidth: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcode", "width_pixels"),
			"Output width of an active transcode in pixels.",
			[]string{"server", "session_id", "user"},
			nil,
		),

		TranscodeHeight: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "transcode", "height_pixels"),
			"Output height of an active transcode in pixels.",
			[]string{"server", "session_id", "user"},
			nil,
		),

		ss: ss,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ds := []*prometheus.Desc{
		c.Info,
		c.DeviceCount,
		c.MovieCount,
		c.SeriesCount,
		c.EpisodeCount,
		c.UserCount,
		c.FailedLoginCount,
		c.MediaStreamCount,
		c.TranscodingCount,
		c.TranscodingBitRate,
		c.VideoBitRate,
		c.AudioBitRate,
		c.AudioChannels,
		c.Up,
		c.ScrapeDuration,
		c.LibraryItemCount,
		c.ScheduledTaskInfo,
		c.ScheduledTaskLastDuration,
		c.ScheduledTaskLastRunTimestamp,
		c.PluginUpdateAvailable,
		c.SessionInfo,
		c.SessionPositionSec,
		c.SessionRuntimeSec,
		c.TranscodeInfo,
		c.TranscodeBitRate,
		c.TranscodeCompletion,
		c.TranscodeCpuUsage,
		c.TranscodeFramerate,
		c.TranscodeWidth,
		c.TranscodeHeight,
	}
	for _, d := range ds {
		ch <- d
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	s, err := c.ss.Status()
	duration := time.Since(start).Seconds()

	up := 1.0
	if err != nil {
		up = 0.0
		log.Printf("failed collecting emby metrics: %v", err)
	}

	ch <- prometheus.MustNewConstMetric(c.Up, prometheus.GaugeValue, up)
	ch <- prometheus.MustNewConstMetric(c.ScrapeDuration, prometheus.GaugeValue, duration)

	if err != nil {
		ch <- prometheus.NewInvalidMetric(c.Info, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		c.Info,
		prometheus.GaugeValue,
		1,
		s.ServerName, s.Hostname, s.Version,
	)

	ch <- prometheus.MustNewConstMetric(
		c.DeviceCount,
		prometheus.GaugeValue,
		float64(s.DeviceCount),
		s.ServerName,
	)

	ch <- prometheus.MustNewConstMetric(
		c.MovieCount,
		prometheus.GaugeValue,
		float64(s.MovieCount),
		s.ServerName,
	)

	ch <- prometheus.MustNewConstMetric(
		c.SeriesCount,
		prometheus.GaugeValue,
		float64(s.SeriesCount),
		s.ServerName,
	)

	ch <- prometheus.MustNewConstMetric(
		c.EpisodeCount,
		prometheus.GaugeValue,
		float64(s.EpisodeCount),
		s.ServerName,
	)

	ch <- prometheus.MustNewConstMetric(
		c.UserCount,
		prometheus.GaugeValue,
		float64(s.UserCount),
		s.ServerName,
	)

	for k, v := range s.FailedLoginCount {
		ch <- prometheus.MustNewConstMetric(
			c.FailedLoginCount,
			prometheus.GaugeValue,
			float64(v),
			s.ServerName, k,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		c.MediaStreamCount,
		prometheus.GaugeValue,
		float64(s.MediaStreamCount),
		s.ServerName,
	)

	ch <- prometheus.MustNewConstMetric(
		c.TranscodingCount,
		prometheus.GaugeValue,
		float64(s.TranscodingCount),
		s.ServerName,
	)

	for k, v := range s.TranscodingBitRate {
		ch <- prometheus.MustNewConstMetric(
			c.TranscodingBitRate,
			prometheus.GaugeValue,
			float64(v),
			s.ServerName, k,
		)
	}

	for k, v := range s.VideoBitRate {
		ch <- prometheus.MustNewConstMetric(
			c.VideoBitRate,
			prometheus.GaugeValue,
			float64(v),
			s.ServerName, k,
		)
	}

	for k, v := range s.AudioBitRate {
		ch <- prometheus.MustNewConstMetric(
			c.AudioBitRate,
			prometheus.GaugeValue,
			float64(v),
			s.ServerName, k,
		)
	}

	for k, v := range s.AudioChannels {
		ch <- prometheus.MustNewConstMetric(
			c.AudioChannels,
			prometheus.GaugeValue,
			float64(v),
			s.ServerName, k,
		)
	}

	for _, lib := range s.Libraries {
		emit := func(typ string, count int) {
			if count == 0 {
				return
			}
			ch <- prometheus.MustNewConstMetric(
				c.LibraryItemCount,
				prometheus.GaugeValue,
				float64(count),
				s.ServerName, lib.Name, lib.CollectionType, typ,
			)
		}
		emit("movie", lib.Counts.MovieCount)
		emit("series", lib.Counts.SeriesCount)
		emit("episode", lib.Counts.EpisodeCount)
		emit("song", lib.Counts.SongCount)
		emit("album", lib.Counts.AlbumCount)
		emit("artist", lib.Counts.ArtistCount)
		emit("music_video", lib.Counts.MusicVideoCount)
		emit("trailer", lib.Counts.TrailerCount)
		emit("box_set", lib.Counts.BoxSetCount)
		emit("book", lib.Counts.BookCount)
		emit("game", lib.Counts.GameCount)
		emit("game_system", lib.Counts.GameSystemCount)
		emit("program", lib.Counts.ProgramCount)
		emit("item", lib.Counts.ItemCount)
	}

	for _, t := range s.ScheduledTasks {
		ch <- prometheus.MustNewConstMetric(
			c.ScheduledTaskInfo,
			prometheus.GaugeValue,
			1,
			s.ServerName, t.Name, t.Category, t.State, t.LastStatus,
		)
		if t.HasResult {
			ch <- prometheus.MustNewConstMetric(
				c.ScheduledTaskLastDuration,
				prometheus.GaugeValue,
				t.LastDurationSec,
				s.ServerName, t.Name,
			)
			ch <- prometheus.MustNewConstMetric(
				c.ScheduledTaskLastRunTimestamp,
				prometheus.GaugeValue,
				t.LastRunUnix,
				s.ServerName, t.Name,
			)
		}
	}

	for _, p := range s.PluginUpdates {
		ch <- prometheus.MustNewConstMetric(
			c.PluginUpdateAvailable,
			prometheus.GaugeValue,
			1,
			s.ServerName, p.Name, p.Version,
		)
	}

	for _, sess := range s.Sessions {
		paused := "false"
		if sess.IsPaused {
			paused = "true"
		}
		ch <- prometheus.MustNewConstMetric(
			c.SessionInfo,
			prometheus.GaugeValue,
			1,
			s.ServerName, sess.SessionId, sess.UserName, sess.Client,
			sess.DeviceName, sess.DeviceType, sess.AppVersion,
			sess.MediaType, sess.PlayMethod, paused, sess.ItemType,
		)
		ch <- prometheus.MustNewConstMetric(
			c.SessionPositionSec,
			prometheus.GaugeValue,
			sess.PositionSec,
			s.ServerName, sess.SessionId, sess.UserName,
		)
		ch <- prometheus.MustNewConstMetric(
			c.SessionRuntimeSec,
			prometheus.GaugeValue,
			sess.RuntimeSec,
			s.ServerName, sess.SessionId, sess.UserName,
		)

		if t := sess.Transcoding; t != nil {
			ch <- prometheus.MustNewConstMetric(
				c.TranscodeInfo,
				prometheus.GaugeValue,
				1,
				s.ServerName, sess.SessionId, sess.UserName,
				t.HwDecoder, t.HwEncoder, t.VideoCodec, t.AudioCodec,
				t.Container, t.Reasons,
			)
			ch <- prometheus.MustNewConstMetric(
				c.TranscodeBitRate, prometheus.GaugeValue, float64(t.Bitrate),
				s.ServerName, sess.SessionId, sess.UserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.TranscodeCompletion, prometheus.GaugeValue, t.Completion,
				s.ServerName, sess.SessionId, sess.UserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.TranscodeCpuUsage, prometheus.GaugeValue, t.CpuUsage,
				s.ServerName, sess.SessionId, sess.UserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.TranscodeFramerate, prometheus.GaugeValue, t.Framerate,
				s.ServerName, sess.SessionId, sess.UserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.TranscodeWidth, prometheus.GaugeValue, float64(t.Width),
				s.ServerName, sess.SessionId, sess.UserName,
			)
			ch <- prometheus.MustNewConstMetric(
				c.TranscodeHeight, prometheus.GaugeValue, float64(t.Height),
				s.ServerName, sess.SessionId, sess.UserName,
			)
		}
	}
}
