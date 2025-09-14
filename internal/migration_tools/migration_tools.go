package migration_tools

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"errors"
	"strconv"

	"github.com/karimNafiz/ChatApplication_WebRTC/internal/utils"
)

//TODO: make sure you make zero padding thing dynamic, rn I don't want to do it
/*
	need to do tests on this
	!!! I forgot to check there could be duplicate migration files
*/

func verify_directory(extension string, directory_path string, seq_num_regex string) error {
	re, err := regexp.Compile(seq_num_regex)
	if err != nil {
		return fmt.Errorf("invalid sequence regex: %w", errInternal)
	}

	upFiles := make(map[string]struct{})
	downFiles := make(map[string]struct{})

	extension = strings.TrimSpace(extension)
	if extension == "" {
		return fmt.Errorf("extension can not be empty: %w", errInternal)
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
		return fmt.Errorf("the number of up (%d) and down (%d) files are not symmetric %w", len(upFiles), len(downFiles), errUpAndDownFilesMismatch)
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
		return fmt.Errorf("missing sequences - up: %v, down: %v %w", missingUp, missingDown, errMissingMigrationSequence)
	}

	return nil
}

/*
	the reason we have two functions ReorderMigrationFiles instead of one is because, we don't want the user
	of this function to provide the seq_num_regex, it should be provided internally

*/

func ReorderMigrationFiles(seq_num int, file_name string, extension string, directory_path string) error {
	seq_num_regex := "^(\\d+)"

	return reorder_migration_files(seq_num, file_name, extension, directory_path, seq_num_regex)
}

func reorder_migration_files(seq_num int, file_name string, extension string, directory_path string, regex_string string) error {
	// Verify first (caller asked for it)
	if err := verify_directory(extension, directory_path, regex_string); err != nil {
		if errors.Is(err, errInternal) {
			return fmt.Errorf("sorry internal error")
		}
		return err
	}

	/*
		repeated code

	*/

	// Normalize extension (".sql", ".txt", ...)
	extension = strings.TrimSpace(extension)
	if extension != "" && extension[0] != '.' {
		extension = "." + extension
	}
	upExt := ".up" + extension
	downExt := ".down" + extension

	re := regexp.MustCompile(regex_string)

	entries, _ := os.ReadDir(directory_path)

	// seq -> filename
	upBySeq := make(map[int]string)
	downBySeq := make(map[int]string)

	seqWidth := 0

	// Collect filenames and infer width
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, extension) {
			continue
		}
		switch {
		case strings.HasSuffix(name, upExt):
			seqStr := re.FindString(name)
			seq, _ := strconv.Atoi(seqStr)
			upBySeq[seq] = name
			if seqWidth == 0 && seqStr != "" {
				seqWidth = len(seqStr)
			}
		case strings.HasSuffix(name, downExt):
			seqStr := re.FindString(name)
			seq, _ := strconv.Atoi(seqStr)
			downBySeq[seq] = name
			if seqWidth == 0 && seqStr != "" {
				seqWidth = len(seqStr)
			}
		}
	}

	// With verify_directory guarantees, both maps have same size and 1..N
	/*
		the sequence number starts from 1 (important)
	*/
	N := len(upBySeq)
	if seq_num < 1 {
		seq_num = 1
	}
	// If inserting at the end (N+1), nothing to shift.
	if seq_num > N+1 {
		seq_num = N + 1
	}

	/*
		highley unlikely
	*/
	if seqWidth == 0 {
		seqWidth = 1
	}

	// Helper to compute new filename by replacing the leading digits
	replaceSeq := func(oldName string, newSeq int) string {
		newHead := fmt.Sprintf("%0*d", seqWidth, newSeq)
		return re.ReplaceAllString(oldName, newHead)
	}

	// Shift files DOWNWARD (N..seq_num) to avoid collisions
	/*
		need to fix this part
	*/
	for i := N; i >= seq_num; i-- {
		// up file
		if oldUp, ok := upBySeq[i]; ok {
			newUp := replaceSeq(oldUp, i+1)
			oldPath := filepath.Join(directory_path, oldUp)
			newPath := filepath.Join(directory_path, newUp)
			_ = os.Rename(oldPath, newPath)
			// upBySeq[i+1] = newUp
			// delete(upBySeq, i)
		}
		// down file
		if oldDown, ok := downBySeq[i]; ok {
			newDown := replaceSeq(oldDown, i+1)
			oldPath := filepath.Join(directory_path, oldDown)
			newPath := filepath.Join(directory_path, newDown)
			_ = os.Rename(oldPath, newPath)
			// downBySeq[i+1] = newDown
			// delete(downBySeq, i)
		}
	}

	/*

		right now we are

	*/

	// After shifting (N..seq_num), create the new pair in the opened slot:
	if strings.TrimSpace(file_name) != "" {
		// If seqWidth wasn’t inferred (e.g., empty dir then inserting at 1), pick a default.
		if seqWidth == 0 {
			// You can choose your default; 6 is common for migrations.
			// this is hard coded must change
			seqWidth = 6
		}

		if _, _, err := createMigrationPair(directory_path, seq_num, seqWidth, file_name, extension); err != nil {
			return fmt.Errorf("failed to create new migration pair at seq %d: %w", seq_num, err)
		}
	}

	return nil
}
