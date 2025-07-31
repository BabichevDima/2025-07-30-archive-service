package service

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BabichevDima/2025-07-30-archive-service/internal/models"
	"github.com/BabichevDima/2025-07-30-archive-service/internal/repository"
)

type ArchiveService interface {
	CreateArchive(taskID string, urls []string) error
}

func NewArchiveServiceImpl(repo *repository.TaskRepository, storagePath string) *ArchiveServiceImpl {
	return &ArchiveServiceImpl{
		repo:        repo,
		storagePath: storagePath,
	}
}

type ArchiveServiceImpl struct {
	repo        *repository.TaskRepository
	storagePath string
}

func (s *ArchiveServiceImpl) CreateArchive(taskID string, urls []string) error {
	tmpDir := filepath.Join(s.storagePath, "tmp", taskID)
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	var downloadedFiles []string
	var errors []string

	for _, url := range urls {
		fileName := filepath.Base(url)
		filePath := filepath.Join(tmpDir, fileName)

		resp, err := http.Get(url)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", url, err))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errors = append(errors, fmt.Sprintf("%s: server returned %d", url, resp.StatusCode))
			continue
		}

		outFile, err := os.Create(filePath)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to create file - %v", url, err))
			continue
		}

		_, err = io.Copy(outFile, resp.Body)
		defer outFile.Close()

		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: failed to save content - %v", url, err))
			continue
		}

		downloadedFiles = append(downloadedFiles, filePath)

		log.Printf("Added %d files to archive, %d errors", len(downloadedFiles), len(errors))
	}

	zipPath := filepath.Join(s.storagePath, fmt.Sprintf("%s.zip", taskID))
	if err := s.createZipArchive(zipPath, downloadedFiles); err != nil {
		return fmt.Errorf("zip creation failed: %w", err)
	}

	return s.repo.UpdateTask(taskID, zipPath, models.StatusCompleted, errors)
}

func (s *ArchiveServiceImpl) createZipArchive(zipPath string, files []string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err := s.addFileToZip(zipWriter, file); err != nil {
			return fmt.Errorf("failed to add file %s to zip: %w", file, err)
		}
	}

	return nil
}

func (s *ArchiveServiceImpl) addFileToZip(zipWriter *zip.Writer, filePath string) error {
	fileToZip, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	return err
}