package migration_tools

import (
	"errors"
)

var (
	errInternal error = errors.New("internal error")
	/*
		when a migration file is missing the sequence prefix so if a migration file has not sequence number in its name
	*/
	errMissingSequencePrefix  error = errors.New("migration file is missing sequence prefix")
	errUpAndDownFilesMismatch error = errors.New("up and down files are not symmetric")

	/*
		when a sequence is missing, so for example when you might have sequence 1 to 3 and have 5 to 6, we have
		the sequence four missing
	*/
	errMissingMigrationSequence error = errors.New("missing migration sequence")
)
