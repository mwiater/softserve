// paths.go
package softserve

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ConvertPath attempts to convert a Unix-like path (e.g., /c/foo/bar)
// used in Git Bash/MinGW environments on Windows to a native Windows path (e.g., C:\foo\bar).
// On non-Windows systems, or if the path doesn't match the expected Git Bash pattern,
// it returns the original path unchanged.
func ConvertPath(path string) string {
	if runtime.GOOS == "windows" {
		// Check if it starts with /drive_letter/ (e.g., /c/ or /d/)
		if len(path) >= 3 && path[0] == '/' && (path[1] >= 'a' && path[1] <= 'z' || path[1] >= 'A' && path[1] <= 'Z') && path[2] == '/' {
			driveLetter := strings.ToUpper(string(path[1]))
			// Replace /c/ with C:\ and then convert remaining / to \
			return driveLetter + ":" + strings.ReplaceAll(path[2:], "/", "\\")
		}
		// Handle UNC paths like //server/share.
		// This assumes Git Bash might pass them as //server/share, which is already close to Windows native.
		if strings.HasPrefix(path, "//") && len(path) > 2 && path[2] != '/' {
			return strings.ReplaceAll(path, "/", "\\")
		}
	}
	// Return original path if not on Windows or if no specific Git Bash conversion pattern is matched.
	return path
}

// EnsureAbsoluteAndExists checks if the given path is an absolute path
// and if it already exists as a directory. It returns an error if
// either condition is not met.
// This function now expects a path that has already been converted if necessary.
func EnsureAbsoluteAndExists(path string) error {
	// 1. Check if the path is absolute
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be an absolute path: got \"%s\"", path)
	}

	// 2. Check if the path exists and is a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Path does not exist, which is the desired error condition
			return fmt.Errorf("path does not exist: %s", path)
		}
		// Some other error occurred while trying to stat the path (e.g., permissions)
		return fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	// 3. Check if the existing path is a directory
	if !fileInfo.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	// If all checks pass, the path is an absolute, existing directory
	return nil
}
