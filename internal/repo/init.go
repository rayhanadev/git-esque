package repo

import (
	"fmt"
	"os"
	"path/filepath"
)

func Init() error {
    err := os.Mkdir(".git", 0755)
    if err != nil {
        return err
    }

    dirs := []string{
        "objects",
        "refs/heads",
        "refs/tags",
    }

    for _, dir := range dirs {
        err = os.MkdirAll(filepath.Join(".git", dir), 0755)
        if err != nil {
            return err
        }
    }

    headFile, err := os.Create(filepath.Join(".git", "HEAD"))
    if err != nil {
        return err
    }
    defer headFile.Close()

    _, err = headFile.WriteString("ref: refs/heads/master\n")
    if err != nil {
        return err
    }

    fmt.Println("Initialized empty Git repository in .git/")
    return nil
}
