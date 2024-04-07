package jobstatus

// JobStatus represents different job statuses.
type JobStatus int

// Constants for different job statuses.
const (
	Default JobStatus = iota
	ScanInitiated
	ScanStarted
	CatzCompleted
	ZapCompleted
	ModulesFinished
	ReportFinished
	FilesCopiedToRemote

	ScanCompleted             = 10
	ScanFailed                = 240
	InvalidOpenAPIFile        = 241
	InputFilesNotPresent      = 242
	UserIDJobIDFolderNotExist = 243
)

// EnumMap maps JobStatus constants to their string representations.
var EnumMap = map[JobStatus]string{
	Default:                   "Yet to Start",
	ScanInitiated:             "Initiating Scan",
	ScanStarted:               "Scan Started",
	CatzCompleted:             "Scan Started",
	ZapCompleted:              "Scan Started",
	ModulesFinished:           "Scan Started",
	ReportFinished:            "Scan Started",
	FilesCopiedToRemote:       "Scan Started",
	ScanCompleted:             "Scan Completed",
	ScanFailed:                "Unable To Scan",
	InvalidOpenAPIFile:        "Invalid OpenAPI File",
	InputFilesNotPresent:      "Input File Error",
	UserIDJobIDFolderNotExist: "Unable To Scan",
}

// GetString returns the string representation of a JobStatus value.
func (j JobStatus) GetText() string {
	if str, ok := EnumMap[j]; ok {
		return str
	}
	return "Invalid Job Status"
}
