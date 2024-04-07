"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

from app.db.connections.maria_db import MariaDatabaseWrapper
from app import core_app


class MainDatabase(MariaDatabaseWrapper):
    """
    * Main Database for the app

    Args:
        Database (_type_): _description_
    """

    def __init__(self):

        super().__init__(
            db_config=core_app.settings.db_config.main, handle_exception=True
        )
