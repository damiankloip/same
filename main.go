package main

import (
  "os"
  "fmt"
  "io"
  "path/filepath"
  "crypto/sha256"
  "github.com/codegangsta/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "Same"
  app.Version = "0.0.1"
  app.Usage = "Find duplicate files in a directory"
  app.UsageText = "same [PATH]"

  app.Action = func(c *cli.Context) {
    base_path := c.Args().First()

    // If the base path is still empty, get the cwd.
    if base_path == "" {
      var err error
      // Default root to cwd.
      base_path, err = os.Getwd()
      check_error(err)
    }

    find_dupes(base_path)
  }

  app.Run(os.Args)
}

func find_dupes(base_path string) {
  hash_map := make(map[string][]string)

  err := filepath.Walk(base_path, func (path string, fileInfo os.FileInfo, file_err error) error {
    if fileInfo.IsDir() {
      return nil;
    }

    check_error(file_err)

    sha, err := compute_sha(path)
    sha_string := string(sha)
    check_error(err)

    hash_map[sha_string] = append(hash_map[sha_string], path)

    return err
  })

  check_error(err)

  for _, matches := range hash_map {
    if (len(matches) > 1) {
      fmt.Println("----")

      for _, filepath := range matches {
        fmt.Println(filepath)
      }
    }
  }
}

func compute_sha(filePath string) ([]byte, error) {
  var result []byte
  file, err := os.Open(filePath)
  defer file.Close()
  check_error(err)

  hash := sha256.New()
  if _, err := io.Copy(hash, file); err != nil {
    return result, err
  }

  return hash.Sum(result), nil
}

// Checks and handles errors.
func check_error(err error) {
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}
