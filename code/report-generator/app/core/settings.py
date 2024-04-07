"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

from dataclasses import dataclass


@dataclass
class MariaDatabaseConfig:
    """
    * Config class format for mariadb
    """

    host: str
    db_name: str
    user: str
    password: str


@dataclass
class AppDatabaseConfig:
    """
    * To store app level db config
    """

    main: MariaDatabaseConfig


@dataclass
class AppSettings:  # pylint: disable=too-many-instance-attributes
    """
    * Class to store app level settings
    """

    db_config: AppDatabaseConfig

    # Directories Reference
    config_dir: str
    app_data_dir: str
    logs_dir: str
    local_temp_dir: str
    user_dir: str
    main_config_path: str

    remote_work_dir: str
