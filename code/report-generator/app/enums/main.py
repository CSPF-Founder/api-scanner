"""
Copyright (c) 2023 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""


class JobStatus:
    DEFAULT = 0
    SCAN_INITIATED = 1
    SCAN_STARTED = 2
    CATZ_COMPLETED = 3
    ZAP_COMPLETED = 4
    MODULES_FINISHED = 5
    REPORT_FINISHED = 6
    FILES_COPIED_TO_REMOTE = 7
    SCAN_FINISHED = 10

    SCAN_FAILED = 240
    INVALID_OPENAPI_FILE = 241
    INPUT_FILES_NOT_PRESENT = 242
    USERID_JOBID_FOLDER_NOT_EXIST = 243
