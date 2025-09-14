package migration_tools

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func get_migration_sequence_frm_filename(re *regexp.Regexp, file_name string) (string, error) {
	match := re.FindString(file_name)

	if match == "" {
		return "", errMissingSequencePrefix
	}

	return match, nil
}

func zeroPad(n, width int) string {
	// %0*d pads n with leading zeros to the given width
	return fmt.Sprintf("%0*d", width, n)
}

// createMigrationPair creates NNNNNN_name.up.ext and NNNNNN_name.down.ext in directory_path.
// It uses O_EXCL so it will NOT overwrite if the files already exist.
func createMigrationPair(directory_path string, seq int, seqWidth int, fileName string, extension string) (upPath, downPath string, err error) {
	// ensure dir exists
	if err = os.MkdirAll(directory_path, 0o755); err != nil {
		return "", "", err
	}

	// normalize extension: ".sql"
	extension = strings.TrimSpace(extension)
	if extension == "" || extension[0] != '.' {
		extension = "." + extension
	}
	upExt := ".up" + extension
	downExt := ".down" + extension

	// slugify the base name: "Add Index" -> "add_index"
	base := slugify(fileName)
	if base == "" {
		base = "migration"
	}

	seqStr := fmt.Sprintf("%0*d", seqWidth, seq)
	upName := seqStr + "_" + base + upExt
	downName := seqStr + "_" + base + downExt

	upPath = filepath.Join(directory_path, upName)
	downPath = filepath.Join(directory_path, downName)

	// choose file contents
	var upBody, downBody []byte
	if extension == ".sql" {
		upBody = []byte(fmt.Sprintf("-- %s up\n\n", base))
		downBody = []byte(fmt.Sprintf("-- %s down\n\n", base))
	} else {
		upBody = []byte{}
		downBody = []byte{}
	}

	// create atomically; if second fails, clean up the first
	if err = createFileAtomic(upPath, upBody); err != nil {
		return "", "", err
	}
	if err = createFileAtomic(downPath, downBody); err != nil {
		_ = os.Remove(upPath)
		return "", "", err
	}

	return upPath, downPath, nil
}

func createFileAtomic(path string, body []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(body); err != nil {
		_ = f.Close()
		_ = os.Remove(path)
		return err
	}
	return f.Close()
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	spaces := regexp.MustCompile(`\s+`)
	s = spaces.ReplaceAllString(s, "_")
	invalid := regexp.MustCompile(`[^a-z0-9_-]`)
	s = invalid.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-_")
	return s
}
