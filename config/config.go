package config

import (
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type GlobalConfig struct {
	MinPeers               int
	Port                   uint16
	TrackerConnectTimeout  time.Duration
	PeerHandshakeTimeout   time.Duration
	PieceMessageTimeout    time.Duration
	DefaultTrackerInterval uint32
	DownloadDirectory      string
	DataDirectory		   string
	MaxTrackerConnections  int
	MaxFailedRetries       int
	PeerId                 [20]byte
}

var Config GlobalConfig

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	defaultDownloadDir := filepath.Join(home, "Downloads")

	Config = GlobalConfig{
		MinPeers:               10,
		Port:                   6881,
		TrackerConnectTimeout:  5 * time.Second,
		PeerHandshakeTimeout:   20 * time.Second,
		PieceMessageTimeout:    30 * time.Second,
		DefaultTrackerInterval: 1800, // 20 minutes
		DownloadDirectory:      defaultDownloadDir,
		DataDirectory: 			getDataDir(),
		MaxTrackerConnections:  20,
		MaxFailedRetries:       3,
	}
}

func getDataDir() string {
	var baseDir string

	if runtime.GOOS == "linux" {
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			home, _ := os.UserHomeDir()
			baseDir = filepath.Join(home, ".local", "share")
		}
	} else {
		baseDir, _ = os.UserConfigDir()
	}

	dataDir := filepath.Join(baseDir, "clover")
	_ = os.MkdirAll(dataDir, 0755)

	return dataDir
}