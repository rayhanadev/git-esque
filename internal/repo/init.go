package repo

import (
	"os"
)

func Init() error {
    err := os.Mkdir(".git", 0755)
    if err != nil {
        return err
    }

    // Create necessary directories and files
    dirs := []string{
        "objects",
        "refs/heads",
        "refs/tags",
    }

    for _, dir := range dirs {
        err = os.MkdirAll(".git/"+dir, 0755)
        if err != nil {
            return err
        }
    }

    // Create HEAD file
    headFile, err := os.Create(".git/HEAD")
    if err != nil {
        return err
    }
    defer headFile.Close()

    _, err = headFile.WriteString("ref: refs/heads/master\n")
    if err != nil {
        return err
    }

    return nil
}
