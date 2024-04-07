"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

import os
import logging
from logging.handlers import RotatingFileHandler
from dataclasses import dataclass


@dataclass
class LoggerConfig:  # pylint: disable=too-many-instance-attributes
    """
    * configuration for the logger
    """
    log_path: str
    log_name: str
    log_level: int = logging.INFO
    max_verbose_level: int = 1
    log_size: int = 5 * 1024 * 1024
    backup_count: int = 10
    log_label: str = "Default"


class Logger:
    """
    * Logger class
    """
    LOG_FORMATTER = logging.Formatter(
       "%(levelname)s-%(asctime)s :%(message)s\n---------", "%Y-%m-%d %H:%M:%S"
    )

    def __init__(
        self,
        logger_config: LoggerConfig
    ):
        self.logger_config = logger_config

        if self.logger_config.log_name:
            self.logger_config.log_label = self.logger_config.log_name

        if not os.path.exists(os.path.dirname(self.logger_config.log_path)):
            # If directory not exists, create dir.
            os.makedirs(os.path.dirname(self.logger_config.log_path))

        self.logger_object = self.get_log_object()

    def get_log_object(self):
        """
        * get logging object

        Returns:
            _type_: _description_
        """
        logger_object = logging.getLogger(self.logger_config.log_label)
        # log_handler = logging.FileHandler(filename=self.log_path)
        log_handler = RotatingFileHandler(
            filename=self.logger_config.log_path,
            maxBytes=self.logger_config.log_size,
            backupCount=self.logger_config.backup_count
        )
        log_handler.setFormatter(self.LOG_FORMATTER)
        logger_object.addHandler(log_handler)
        logger_object.setLevel(self.logger_config.log_level)
        return logger_object

    def error(self, message, verbose_level=1, label=''):
        """
        * Function to log error

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        if verbose_level > self.logger_config.max_verbose_level:
            return

        if label:
            self.logger_object.error(label + " - " + message)
        else:
            self.logger_object.error(message)

    def exception(self, exception):
        """
        * Function to log exception

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        self.logger_object.exception(exception)

    def info(self, message, verbose_level=1, label=''):
        """
        * Function to log info

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        if verbose_level > self.logger_config.max_verbose_level:
            return

        if label:
            self.logger_object.info(label + " - " + message)
        else:
            self.logger_object.info(message)

    def warning(self, message, verbose_level=1, label=''):
        """
        * Function to log warning

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        if verbose_level > self.logger_config.max_verbose_level:
            return

        if label:
            self.logger_object.warning(label + " - " + message)
        else:
            self.logger_object.warning(message)


class LoggerWithInfoAndErrorFiles:
    """
    * Logger class to log info in separate file and error in separate file
    """

    def __init__(self, info_log_path, exception_log_path, log_name, max_verbose_level):
        self.scan_log_path = info_log_path
        self.exception_log_path = exception_log_path

        self.info_logger = Logger(LoggerConfig(
            log_path=info_log_path,
            log_name=log_name + "_info",
            log_level=logging.INFO,
            max_verbose_level=max_verbose_level
        ))

        self.error_logger = Logger(LoggerConfig(
            log_path=exception_log_path,
            log_name=log_name + "_error",
            log_level=logging.DEBUG,
            max_verbose_level=max_verbose_level
        ))

    def error(self, message, verbose_level=1, label=''):
        """
        * Function to log error

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        self.error_logger.error(message, verbose_level, label)

    def exception(self, exception):
        """
        * Function to log exception

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        self.error_logger.exception(exception)

    def info(self, message, verbose_level=1, label=''):
        """
        * Function to log info

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        self.info_logger.info(message, verbose_level, label)

    def warning(self, message, verbose_level=1, label=''):
        """
        * Function to log warning

        Args:
            message (_type_): _description_
            verbose_level (int, optional): _description_. Defaults to 1.
            label (str, optional): _description_. Defaults to ''.
        """
        self.error_logger.warning(message, verbose_level, label)


class AppLogger(LoggerWithInfoAndErrorFiles):
    """
    * App logger class

    Args:
        LoggerWithInfoAndErrorFiles (_type_): _description_
    """

    def __init__(
        self,
        info_log_path,
        exception_log_path,
        log_name="app_logger",
        max_verbose_level=1
    ):
        super().__init__(
            info_log_path,
            exception_log_path,
            log_name,
            max_verbose_level
        )


class ScanLogger(LoggerWithInfoAndErrorFiles):
    """
    * Scan logger class

    Args:
        LoggerWithInfoAndErrorFiles (_type_): _description_
    """

    def __init__(
        self,
        info_log_path,
        exception_log_path,
        log_name="scan_logger",
        max_verbose_level=1
    ):
        super().__init__(
            info_log_path,
            exception_log_path,
            log_name,
            max_verbose_level
        )
