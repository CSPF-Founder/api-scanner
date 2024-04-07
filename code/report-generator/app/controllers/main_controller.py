import configparser
import os
import logging

from app.db.models.jobs import JobModel
from app.enums.main import JobStatus
from app.modules.reporter import Reporter
from app.utils.custom_error import AppError
from app.validators.openapi import validate_openapi_file
from app.modules.file_sync import (
    copy_to_remote_only_log,
)


class MainController:
    def __init__(self, job_id: int, user_id: int):
        self.job = JobModel(id=job_id, user_id=user_id)

    def handle_report_failure(self, error_message, logger):
        if logger:
            logger.error(error_message)
            try:
                copy_to_remote_only_log(job=self.job)
            except Exception:
                pass

        json_ouptut = {
            "success": False,
            "error": error_message,
        }
        print(json_ouptut)
        exit(1)

    def get_auth_headers(self, auth_header_file, logger):
        """
        * Check if authfile has headers
        """
        auth_config = configparser.ConfigParser()
        auth_config.read(auth_header_file)
        list_of_values = auth_config.items("AUTH_HEADERS")
        number_of_headers = len(list_of_values)
        if number_of_headers == 1:
            return list_of_values[0]
        if number_of_headers == 0:
            logger.info("No Security Headers")
            return False
        if number_of_headers > 1:
            logger.error("Multiple auth headers. Cannot process")
            return False

    def get_logger(self, log_file_path) -> logging.Logger:
        logger = logging.getLogger(f"report_logger_{self.job.id}")
        log_handler = logging.FileHandler(
            filename=log_file_path,
        )
        log_handler.setLevel(logging.INFO)
        log_handler.setFormatter(logging.Formatter("%(asctime)s %(message)s"))
        logger.addHandler(log_handler)
        return logger

    def run(self):
        is_report_generated = False
        logger = None
        try:
            scan_dir = self.job.get_local_work_dir()
            logs_file = os.path.abspath(scan_dir + "report.log")
            auth_header_file = os.path.abspath(scan_dir + "auth_headers.conf")
            openapi_file = os.path.abspath(scan_dir + "openapi.yaml")

            # Create directories
            if not os.path.exists(self.job.get_local_work_dir()):
                os.makedirs(self.job.get_local_work_dir())

            logger = self.get_logger(logs_file)

            if not os.path.exists(openapi_file) or not os.path.exists(auth_header_file):
                raise AppError(
                    message="Input files not present in expected path",
                    status=JobStatus.INPUT_FILES_NOT_PRESENT,
                )

            auth_headers_data = self.get_auth_headers(auth_header_file, logger)

            server_url = None
            try:
                is_valid_openapi_result, server_url = validate_openapi_file(
                    openapi_file
                )
            except Exception as e:
                logger.error("Validation Exception" + str(e))
                is_valid_openapi_result = False

            if not is_valid_openapi_result:
                raise AppError(
                    message="Invalid OpenAPI File",
                    status=JobStatus.INVALID_OPENAPI_FILE,
                )

            reporter = Reporter(job=self.job)
            is_report_generated = reporter.run(
                server_url=server_url,
                auth_headers_data=auth_headers_data,
                logger=logger,
            )

        except AppError as error:
            self.handle_report_failure(error_message=error.message, logger=logger)
        finally:
            if is_report_generated:
                json_ouptut = {
                    "success": True,
                }
                print(json_ouptut)
            else:
                self.handle_report_failure(
                    error_message="Report generation failed",
                    logger=logger,
                )
