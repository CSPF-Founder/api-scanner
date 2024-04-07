"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

from configparser import ConfigParser
import os
import time

from app.core.config import AppConfig
from app.core.settings import AppSettings
from app.core.settings import AppDatabaseConfig, MariaDatabaseConfig
from app.core.logger import AppLogger


class App:
    """Core app object

    Args:
        object (_type_): _description_
    """

    def __init__(self):
        self.config = AppConfig()

        self.settings: AppSettings = self.load_config()

        self.logger = AppLogger(
            os.path.join(self.settings.logs_dir, "app_info.log"),
            os.path.join(self.settings.logs_dir, "app_error.log"),
            log_name="app_logger",
            # max_verbose_level=self.config.log_verbose_level
        )

    def load_config(self) -> AppSettings:
        """
        * Load config_parser_object from app.conf and put in dictionary
        """
        # interpolation is disabled
        config_parser_object = ConfigParser(interpolation=None)
        config_parser_object.read(self.get_main_config_path())

        return self.load_base_config(config_parser_object)

    def get_main_config_dir(self):
        """
        Calculate/format main config directory

        Returns:
            _type_: _description_
        """
        return os.path.join(self.config.app_path, "config")

    def get_main_config_path(self):
        """
        Calculate/format main config directory

        Returns:
            _type_: _description_
        """
        return os.path.join(self.get_main_config_dir(), "app.conf")

    def load_base_config(self, config_obj) -> AppSettings:
        """
        * Load basic configuration from the app.conf

        Args:
            config_obj (_type_): _description_
        """

        # Directorie path reference
        app_data_dir = os.path.join(self.config.app_path, "app_data")
        # local_temp_dir = os.path.join(self.config.app_path, "local_temp")
        local_temp_dir = config_obj.get("MAIN", "local_temp_dir")
        if not os.path.exists(local_temp_dir):
            os.makedirs(local_temp_dir)

        remote_work_dir = config_obj.get("MAIN", "remote_work_dir")

        logs_dir = os.path.join(self.config.app_path, "logs")
        user_dir = os.path.join(self.config.app_path, "user_data")

        # Initialize Settings
        return AppSettings(
            db_config=self.load_db_config(config_obj),
            app_data_dir=app_data_dir,
            local_temp_dir=local_temp_dir,
            config_dir=self.get_main_config_dir(),
            main_config_path=self.get_main_config_path(),
            logs_dir=logs_dir,
            user_dir=user_dir,
            remote_work_dir=remote_work_dir,
        )

    def load_db_config(self, config_obj) -> AppDatabaseConfig:
        """
        * Load Database configs

        Args:
            config_obj (_type_): _description_

        Returns:
            AppDatabaseConfig: _description_
        """
        main_db_config = MariaDatabaseConfig(
            host=config_obj.get("MAIN_DATABASE", "host"),
            db_name=config_obj.get("MAIN_DATABASE", "db_name"),
            user=config_obj.get("MAIN_DATABASE", "user"),
            password=config_obj.get("MAIN_DATABASE", "password"),
        )

        return AppDatabaseConfig(
            main=main_db_config,
        )


def create_app():
    """
    * Function to create core app object
    """
    print("Initiating App")

    # Set Timezone
    os.environ["TZ"] = "Asia/Calcutta"
    time.tzset()

    app_obj = App()
    return app_obj


core_app = create_app()
