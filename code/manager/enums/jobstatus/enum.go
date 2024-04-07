package jobstatus

// JobStatus represents different job statuses.
type JobStatus int32

// Constants for different job statuses.
const (
	Default             JobStatus = 0
	ScanInitiated       JobStatus = 1
	ScanStarted         JobStatus = 2
	CatzCompleted       JobStatus = 3
	ZapCompleted        JobStatus = 4
	ModulesFinished     JobStatus = 5
	ReportFinished      JobStatus = 6
	FilesCopiedToRemote JobStatus = 7

	ScanCompleted             JobStatus = 10
	ScanFailed                JobStatus = 240
	InvalidOpenAPIFile        JobStatus = 241
	InputFilesNotPresent      JobStatus = 242
	UserIDJobIDFolderNotExist JobStatus = 243
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
