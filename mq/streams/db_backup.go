package streams

import "github.com/nats-io/nats.go"

const (
	StreamDbBackup string = "DbBackup"

	SubjectDbBackupExec = "DBBACKUP.Execute"

	RetentionPolicyDbBackup = nats.WorkQueuePolicy

	DuplicateFilterWindowDbBackup = 0
)

var StreamSubjectsDbBackup = []string{
	SubjectDbBackupExec,
}
