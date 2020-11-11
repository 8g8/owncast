package storageproviders

import (
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/owncast/owncast/config"
	"github.com/owncast/owncast/core/ffmpeg"
	"github.com/owncast/owncast/utils"
)

type LocalStorage struct {
}

// Setup configures this storage provider
func (s *LocalStorage) Setup() error {
	return nil
}

// NOTE: This cleanup timer will have to be disabled to support recordings in the future
// as all HLS segments have to be publicly available on disk to keep a recording of them.

func (s *LocalStorage) Cleanup() {
	ffmpeg.CleanupOldContent(config.PublicHLSStoragePath)
}

// SegmentWritten is called when a single segment of video is written
func (s *LocalStorage) SegmentWritten(localFilePath string) {
	s.Save(localFilePath, 0)
}

// VariantPlaylistWritten is called when a variant hls playlist is written
func (s *LocalStorage) VariantPlaylistWritten(localFilePath string) {
	_, err := s.Save(localFilePath, 0)
	if err != nil {
		log.Errorln(err)
		return
	}
}

// MasterPlaylistWritten is called when the master hls playlist is written
func (s *LocalStorage) MasterPlaylistWritten(localFilePath string) {
	s.Save(localFilePath, 0)
}

// Save will save a local filepath using the storage provider
func (s *LocalStorage) Save(filePath string, retryCount int) (string, error) {
	newPath := ""

	// This is a hack
	if filePath == "hls/stream.m3u8" {
		newPath = filepath.Join(config.PublicHLSStoragePath, filepath.Base(filePath))
	} else {
		newPath = filepath.Join(config.WebRoot, filePath)
	}

	utils.Copy(filePath, newPath)

	return newPath, nil
}
