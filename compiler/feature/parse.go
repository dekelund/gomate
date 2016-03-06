// Package feature provides ...
package feature

import (
	"os"
	"path/filepath"

	. "gomate.io/gomate/global"
)

// New allocates and returns a string slice, containing feature-files
// related to packages provided as argument from caller.
func New(pkgs []string) (features []string) {
	paths := PKGPathToPath(Settings.CWD, Settings.GOSRCPATH, pkgs)

	for _, path := range paths {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			_, file := filepath.Split(path)

			if filepath.Ext(file) == ".feature" {
				features = append(features, path)
			}

			return nil
		})
	}

	return
}
