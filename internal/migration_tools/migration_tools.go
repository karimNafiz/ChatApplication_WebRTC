package migration_tools

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/karimNafiz/ChatApplication_WebRTC/internal/utils"
)

//TODO: make sure you make zero padding thing dynamic, rn I don't want to do it

func verify_directory(extension string, directory_path string, seq_num_regex string) error {
	re, err := regexp.Compile(seq_num_regex)
	if err != nil {
		return fmt.Errorf("invalid sequence regex: %w", err)
	}

	upFiles := make(map[string]struct{})
	downFiles := make(map[string]struct{})

	extension = strings.TrimSpace(extension)
	if extension == "" {
		return fmt.Errorf("extension can not be empty")
	}
	if extension[0] != '.' {
		extension = "." + extension
	}

	entries, err := os.ReadDir(directory_path)
	if err != nil {
		return fmt.Errorf("could not open directory %q: %w", directory_path, err)
	}

	upExt := ".up" + extension
	downExt := ".down" + extension

	// for inferring width from any seen sequence (e.g., "000001" -> width=6)
	seqWidth := 0

	for _, file := range entries {
		name := file.Name()

		// only consider files that end with the given extension
		if !strings.HasSuffix(name, extension) {
			continue
		}

		switch {
		case strings.HasSuffix(name, upExt):
			seq, err := get_migration_sequence_frm_filename(re, name)
			if err != nil {
				return fmt.Errorf("migration file missing sequence number %q: %w", name, err)
			}
			upFiles[seq] = struct{}{}
			// I don't like this part of the code but necessary evil
			if seqWidth == 0 {
				seqWidth = len(seq)
			}
		case strings.HasSuffix(name, downExt):
			seq, err := get_migration_sequence_frm_filename(re, name)
			if err != nil {
				return fmt.Errorf("migration file missing sequence number %q: %w", name, err)
			}
			downFiles[seq] = struct{}{}
			// I don't like this part of the code but necessary evil
			if seqWidth == 0 {
				seqWidth = len(seq)
			}
		}
	}

	// counts must be equal and non-zero (optional non-zero check)
	if !(len(upFiles) == len(downFiles) && len(utils.Intersection(upFiles, downFiles)) == len(upFiles)) {
		return fmt.Errorf("the number of up (%d) and down (%d) files are not symmetric", len(upFiles), len(downFiles))
	}

	if len(upFiles) == 0 {
		// nothing to check further
		return nil
	}

	// Check for contiguous sequences: 1..N (zero-padded)
	N := len(upFiles)

	var missingUp []string
	var missingDown []string

	for i := 1; i <= N; i++ {
		seq := zeroPad(i, seqWidth) // e.g., i=2, width=6 -> "000002"

		if _, ok := upFiles[seq]; !ok {
			missingUp = append(missingUp, seq)
		}
		if _, ok := downFiles[seq]; !ok {
			missingDown = append(missingDown, seq)
		}
	}

	if len(missingUp) > 0 || len(missingDown) > 0 {
		return fmt.Errorf("missing sequences - up: %v, down: %v", missingUp, missingDown)
	}

	return nil
}

func zeroPad(n, width int) string {
	// %0*d pads n with leading zeros to the given width
	return fmt.Sprintf("%0*d", width, n)
}

func reorder_migration_files(seq_num int, file_name string, extension string, directory_path string, seq_num_regex string) error {
	return nil
}

func get_migration_sequence_frm_filename(re *regexp.Regexp, file_name string) (string, error) {
	match := re.FindString(file_name)

	if match == "" {
		return "", fmt.Errorf("there is no sequence number in the file name")
	}

	return match, nil
}
