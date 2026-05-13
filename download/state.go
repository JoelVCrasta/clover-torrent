package download

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/JoelVCrasta/clover/client"
	"github.com/JoelVCrasta/clover/config"
)

var (
	DATA_DIR     = config.Config.DataDirectory
	STATE_DIR    = filepath.Join(DATA_DIR, "state")
	SESSION_FILE = filepath.Join(DATA_DIR, "session.json")
)

type Session struct {
	Infohash    string `json:"info_hash"`
	Name        string `json:"name"`
	Done        int    `json:"done"`
	Total       int    `json:"total"`
	TimeElapsed int64  `json:"time_elapsed"`
	InputPath   string `json:"input_path"`
	OutputPath  string `json:"output_path"`
}

func (dm *DownloadManager) SaveSession() error {
	dir := filepath.Dir(DATA_DIR)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	var sessions []Session
	data, err := os.ReadFile(SESSION_FILE)

	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &sessions); err != nil {
			return fmt.Errorf("Failed to read session: %w", err)
		}
	}

	currInfoHash := fmt.Sprintf("%x", dm.torrent.InfoHash)

	newSession := Session{
		Infohash:    currInfoHash,
		Name:        dm.torrent.Info.Name,
		Done:        dm.stats.Done,
		Total:       dm.stats.Total,
		TimeElapsed: int64(dm.stats.TimeElapsed.Seconds()),
		InputPath:   filepath.Join(STATE_DIR, currInfoHash, currInfoHash+".torrent"),
		OutputPath:  GetOutputBasePath(dm.torrent),
	}

	found := false
	for i, s := range sessions {
		if s.Infohash == currInfoHash {
			sessions[i] = newSession
			found = true
			break
		}
	}

	if !found {
		sessions = append(sessions, newSession)
	}

	updatedData, err := json.MarshalIndent(sessions, "", "	")
	if err != nil {
		return fmt.Errorf("Failed to save session: %w", err)
	}

	return os.WriteFile(SESSION_FILE, updatedData, 0644)
}

func (dm *DownloadManager) LoadState() error {

	return nil
}

func (dm *DownloadManager) SaveState() error {
	currInfoHash := fmt.Sprintf("%x", dm.torrent.InfoHash)

	torrentDir := filepath.Join(STATE_DIR, currInfoHash)
	torrentFile := filepath.Join(torrentDir, currInfoHash+".torrent")
	bitfieldFile := filepath.Join(torrentDir, currInfoHash+".bitfield")

	if err := os.MkdirAll(torrentDir, 0755); err != nil {
		return fmt.Errorf("Failed to create state dir: %w", err)
	}

	if _, err := os.Stat(torrentFile); os.IsNotExist(err) {
		if len(dm.torrent.BencodeByteStream) > 0 {
			if err := os.WriteFile(torrentFile, dm.torrent.BencodeByteStream, 0644); err != nil {
				return fmt.Errorf("Failed to save torrent metadata: %w", err)
			}
		}
	}
	dm.torrent.BencodeByteStream = nil

	bf := make(client.Bitfield, (dm.stats.Total+7)/8)
	for i, done := range dm.downloadedPieces {
		if done {
			bf.Set(i)
		}
	}

	return os.WriteFile(bitfieldFile, bf, 0644)
}
