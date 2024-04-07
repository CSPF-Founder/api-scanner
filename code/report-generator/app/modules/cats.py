import os
import re
import shutil
import json

# Internal dependencies
from app.db.models.jobs import JobModel
from app.enums.main import JobStatus
from app.utils import common_utils
from app.utils.custom_error import AppError


class CatsReporter:
    def __init__(self, job: JobModel):
        self.job = job
        self.output_dir = os.path.join(self.job.get_local_work_dir(), "catsoutput/")
        self.requests_dir = os.path.join(self.job.get_local_work_dir(), "requests/")

    # add catz data to report
    def add_to_report(self, document, logger):
        count = 0

        for root, _, files in os.walk(self.output_dir):
            for file_name in files:
                if not file_name.endswith(".json"):
                    continue

                file_path = os.path.join(root, file_name)
                try:
                    count = self.add_to_report_from_json(
                        count=count,
                        document=document,
                        file_name=file_name,
                        file_path=file_path,
                        logger=logger,
                    )
                except Exception as e:
                    logger.error("Error in adding to report." + str(e))
                    continue
            if count == 0:
                document.add_heading("Fuzz Scan Results:", 1)
                document.add_heading("No issues found", level=2)
                document.add_page_break()

    def add_to_report_from_json(
        self, *, count: int, document, file_name: str, file_path: str, logger
    ) -> int:
        with open(file_path, "r") as f:
            json_data = json.load(f)

        result = json_data["result"]

        if result != "error":
            return count

        # copy request file for customer
        self.copy_requests_file(file_name=file_name, src_file=file_path)

        count = count + 1
        if count == 1:
            document.add_heading("Fuzz Scan Results:", 1)
            document.add_heading(
                "Why do unprocessed exceptions need to be fixed? \t\t:", level=2
            )
            document.add_paragraph(
                "The scanner points out errors such as 500 errors. "
                "A 500 error is the result of an unexpected condition. "
                "This means that the API is not taking into account all "
                "input possibilities. API’s typically should not have any "
                "unhandled exceptions. All types of inputs should be "
                "properly handled or given specific errors."
            )
            document.add_paragraph(
                "All such errors should be handled with specific error"
                " codes or messages. This helps make the API more stable"
                " and reduces the chances of vulnerabilities."
            )
            document.add_page_break()

        if count == 100:
            document.add_heading(
                "Too many errors to show in report."
                " To see the rest of the requests please refer "
                "to error_requests.zip file"
            )
            document.add_page_break()
            return count

        if count > 100:
            return count

        document.add_heading(
            str(count) + "." + "  Tested scenario :   "
            "" + self.clean_cats_from_string(json_data["fuzzer"]),
            level=2,
        )
        document.add_heading("Result \t\t:", level=2)
        document.add_paragraph(self.clean_cats_from_string(json_data["resultReason"]))
        document.add_heading("Replication File \t\t: " + file_name, level=2)

        (curl_string, error) = self.create_curl(json_data["request"])

        if curl_string:
            document.add_heading("Replication CURL \t\t:", level=2)
            document.add_paragraph(curl_string)
        else:
            document.add_heading("Replication \t\t:", level=2)
            if error:
                document.add_paragraph(
                    "Cannot print in document due to "
                    + str(error)
                    + f". Refer file {file_name} in error_requests.zip"
                )
            else:
                document.add_paragraph(f"Refer file {file_name} in error_requests.zip")

        return count

    # Copy error json request file and sanatize catz => qats
    def copy_requests_file(self, file_name, src_file):
        dst_file = os.path.join(self.requests_dir, file_name)
        shutil.copyfile(src_file, dst_file)

        with open(dst_file, "rt") as fin:
            datar = fin.read()
            datar = datar.replace("cats", "qats")

        if not datar:
            return

        with open(dst_file, "wt") as fw:
            fw.write(datar)
            fw.close()

    # Used in report
    def clean_cats_from_string(self, input_value: str):
        pattern = re.compile("cats", re.IGNORECASE)
        return common_utils.replace_unsupported_for_docx(
            common_utils.smart_str(pattern.sub("qats", input_value))
        )

    # Used in report
    def sanitize_word_or_error(self, input_value: str):
        if common_utils.any_unsupported_for_docx(input_value):
            raise AppError(
                "unsupported character for docx", JobStatus.INVALID_OPENAPI_FILE
            )

        pattern = re.compile("cats", re.IGNORECASE)
        return common_utils.replace_unsupported_for_docx(
            common_utils.smart_str(pattern.sub("qats", input_value))
        )

    # Create a CURL request for
    def create_curl(self, json_cats_request) -> tuple:
        curl_string = ""
        try:
            http_method = self.sanitize_word_or_error(json_cats_request["httpMethod"])

            target_url = self.sanitize_word_or_error(json_cats_request["url"])
            curl_string = f"curl -X '{http_method}' --url {target_url} "

            post_data = self.sanitize_word_or_error(json_cats_request["payload"])
            http_method = json_cats_request["httpMethod"]
            header_data = json_cats_request["headers"]

            if http_method != "GET" and (post_data != "{}" and post_data is not None):
                if len(post_data) > 1000:
                    # post_data = '{"STRINGTOOLONG":"Check error_requests.zip file"}'
                    return (None, "large post data")
                curl_string = f"{curl_string} -d '{post_data}'"

            if len(header_data) > 15:
                # curl_string = (
                #     f"{curl_string} -H 'TooManyHeaders:Check error_requests.zip file'"
                # )
                return (None, "too many headers")
            else:
                for header in header_data:
                    header_key = self.sanitize_word_or_error(header["key"])
                    header_value = self.sanitize_word_or_error(header["value"])

                    curl_string = f"{curl_string} -H '{header_key}:{header_value}'"

        except AppError as e:
            return (None, str(e))
        except Exception:
            return (None, None)

        return (curl_string, None)
