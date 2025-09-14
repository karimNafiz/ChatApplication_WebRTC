package migration_tools

import (
	"fmt"
	"regexp"
)

func get_migration_sequence_frm_filename(re *regexp.Regexp, file_name string) (string, error) {
	match := re.FindString(file_name)

	if match == "" {
		return "", fmt.Errorf("there is no sequence number in the file name")
	}

	return match, nil
}

func zeroPad(n, width int) string {
	// %0*d pads n with leading zeros to the given width
	return fmt.Sprintf("%0*d", width, n)
}
