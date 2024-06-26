package repo

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func Clone(repoURL string, directory string) error {
    if !strings.HasSuffix(repoURL, ".git") {
        repoURL += ".git"
    }
    tarballURL := fmt.Sprintf("%s/archive/master.tar.gz", strings.TrimSuffix(repoURL, ".git"))

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

    for {
        header, err := tarReader.Next()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }

        target := filepath.Join(directory, header.Name)
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

    fmt.Printf("Cloned repository to %s\n", directory)
    return nil
}
