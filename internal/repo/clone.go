package repo

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Clone(repoURL, directory, branch string) error {
	if !strings.HasSuffix(repoURL, ".git") {
		repoURL += ".git"
	}

	if branch == "" {
		branch = "master"
	}

	tarballURL := fmt.Sprintf("%s/archive/refs/heads/%s.tar.gz", strings.TrimSuffix(repoURL, ".git"), branch)

	if directory == "" {
		parts := strings.Split(repoURL, "/")
		directory = strings.TrimSuffix(parts[len(parts)-1], ".git")
	}

	err := os.Mkdir(directory, 0755)
	if err != nil {
		return err
	}

	response, err := http.Get(tarballURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	gzReader, err := gzip.NewReader(response.Body)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	var baseDir string

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if baseDir == "" && header.Typeflag != tar.TypeXGlobalHeader {
			baseDir = strings.Split(header.Name, "/")[0]
		}

		relativePath := strings.TrimPrefix(header.Name, baseDir+"/")
		if relativePath == "" {
			continue
		}
		target := filepath.Join(directory, relativePath)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		case tar.TypeXGlobalHeader:
			// Ignore pax_global_header
			continue
		default:
			return fmt.Errorf("unknown type: %b in %s", header.Typeflag, header.Name)
		}
	}

	gitDir := filepath.Join(directory, ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		return err
	}
	subdirs := []string{"objects", "refs/heads", "refs/tags", "objects/pack"}
	for _, subdir := range subdirs {
		err = os.MkdirAll(filepath.Join(gitDir, subdir), 0755)
		if err != nil {
			return err
		}
	}

	headURL := fmt.Sprintf("%s/HEAD", repoURL)
	err = fetchAndWriteFile(headURL, filepath.Join(gitDir, "HEAD"))
	if err != nil {
		return err
	}

	refsURL := fmt.Sprintf("%s/info/refs", repoURL)
	err = fetchAndWriteRefs(refsURL, gitDir)
	if err != nil {
		return err
	}

	err = fetchObjects(repoURL, gitDir)
	if err != nil {
		return err
	}

	fmt.Printf("Cloned repository to %s\n", directory)
	return nil
}

func fetchAndWriteFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func fetchAndWriteRefs(refsURL, gitDir string) error {
	resp, err := http.Get(refsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) == 2 {
			hash := parts[0]
			ref := parts[1]
			refPath := filepath.Join(gitDir, ref)
			err = os.MkdirAll(filepath.Dir(refPath), 0755)
			if err != nil {
				return err
			}
			err = os.WriteFile(refPath, []byte(hash), 0644)
			if err != nil {
				return err
			}
		}
	}

	return scanner.Err()
}

func fetchObjects(repoURL, gitDir string) error {
	packsURL := fmt.Sprintf("%s/objects/info/packs", repoURL)
	resp, err := http.Get(packsURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "P pack-") {
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				packFile := parts[1]
				err = downloadPackFile(repoURL, packFile, gitDir)
				if err != nil {
					return err
				}
			}
		}
	}

	return scanner.Err()
}

func downloadPackFile(repoURL, packFile, gitDir string) error {
	packURL := fmt.Sprintf("%s/objects/pack/%s", repoURL, packFile)
	resp, err := http.Get(packURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	targetPath := filepath.Join(gitDir, "objects", "pack", packFile)
	outFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resp.Body)
	return err
}
