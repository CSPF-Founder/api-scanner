"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""
import os


class AppConfig:
    """
    Application Configuration
    """

    def __init__(self):
        self._scanner_path = None

    @property
    def scanner_path(self):
        if not self._scanner_path:
            current_file_directory = os.path.dirname(
                os.path.realpath(__file__))
            self._scanner_path = os.path.dirname(
                os.path.dirname(current_file_directory))
        return self._scanner_path

    @property
    def app_path(self):
        return os.path.dirname(self.scanner_path)

   