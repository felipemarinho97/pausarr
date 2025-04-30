package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/autobrr/go-qbittorrent"
	"github.com/joho/godotenv"
)

type Config struct {
	QBittorrentURL      string
	QBittorrentUsername string
	QBittorrentPassword string
	MaxRatio            float64
	MaxSeedDuration     time.Duration
	SkipPrivateTrackers bool
	CheckInterval       time.Duration
}

func loadConfig() *Config {
	_ = godotenv.Load()

	maxRatio, _ := strconv.ParseFloat(os.Getenv("MAX_RATIO"), 64)
	if maxRatio == 0 {
		maxRatio = 2.0
	}

	maxSeedHours, _ := strconv.Atoi(os.Getenv("MAX_SEED_HOURS"))
	if maxSeedHours == 0 {
		maxSeedHours = 24 * 30 * 6 // 6 months
	}

	skipPrivate, _ := strconv.ParseBool(os.Getenv("SKIP_PRIVATE_TRACKERS"))

	checkInterval, _ := strconv.Atoi(os.Getenv("CHECK_INTERVAL_MINUTES"))
	if checkInterval == 0 {
		checkInterval = 60 // 1 hour
	}

	log.Printf("Config loaded: MaxRatio=%f, MaxSeedDuration=%d, SkipPrivateTrackers=%t, CheckInterval=%d",
		maxRatio, maxSeedHours, skipPrivate, checkInterval)

	return &Config{
		QBittorrentURL:      os.Getenv("QBITTORRENT_URL"),
		QBittorrentUsername: os.Getenv("QBITTORRENT_USERNAME"),
		QBittorrentPassword: os.Getenv("QBITTORRENT_PASSWORD"),
		MaxRatio:            maxRatio,
		MaxSeedDuration:     time.Duration(maxSeedHours) * time.Hour,
		SkipPrivateTrackers: skipPrivate,
		CheckInterval:       time.Duration(checkInterval) * time.Minute,
	}
}

func main() {
	config := loadConfig()

	client := qbittorrent.NewClient(qbittorrent.Config{
		Host:     config.QBittorrentURL,
		Username: config.QBittorrentUsername,
		Password: config.QBittorrentPassword,
	})

	ctx := context.Background()

	// Test initial connection
	if err := client.LoginCtx(ctx); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	log.Println("Connection established successfully")

	ticker := time.NewTicker(config.CheckInterval)
	defer ticker.Stop()

	for {
		processTorrents(ctx, client, config)
		<-ticker.C
	}
}

func processTorrents(ctx context.Context, client *qbittorrent.Client, config *Config) {
	torrents, err := client.GetTorrents(qbittorrent.TorrentFilterOptions{
		Filter: qbittorrent.TorrentFilterUploading,
	})
	if err != nil {
		log.Printf("Error fetching torrents: %v", err)
		return
	}
	log.Printf("Fetched %d torrents from qbittorrent", len(torrents))

	pauseCount := 0

	for _, torrent := range torrents {
		if shouldPauseTorrent(&torrent, config, client) {
			pauseCount++
			if err := client.PauseCtx(ctx, []string{torrent.Hash}); err != nil {
				log.Printf("Error pausing %s: %v", torrent.Name, err)
				continue
			}

			log.Printf("[PAUSED] %s - Ratio: %.2f, Seed: %s", torrent.Name, torrent.Ratio, time.Duration(torrent.SeedingTime)*time.Second)
		}
	}

	if pauseCount > 0 {
		log.Printf("Paused %d torrents", pauseCount)
	}
}

func shouldPauseTorrent(torrent *qbittorrent.Torrent, config *Config, client *qbittorrent.Client) bool {
	// Ignore already paused torrents
	if torrent.State == qbittorrent.TorrentStatePausedDl || torrent.State == qbittorrent.TorrentStatePausedUp {
		return false
	}

	// Check if it has the "no-pause" tag
	if strings.Contains(torrent.Tags, "no-pause") {
		return false
	}

	// Check private trackers
	if config.SkipPrivateTrackers {
		tp, err := client.GetTorrentProperties(torrent.Hash)
		if err != nil {
			log.Printf("Error fetching properties for torrent %s: %v", torrent.Name, err)
			return false
		}
		if tp.IsPrivate {
			return false
		}
	}

	// Check seed time & ratio
	seedDuration := time.Duration(torrent.SeedingTime) * time.Second
	return torrent.Ratio >= config.MaxRatio || seedDuration >= config.MaxSeedDuration
}
