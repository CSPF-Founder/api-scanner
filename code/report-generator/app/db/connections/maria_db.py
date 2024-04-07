"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""

import mysql.connector as mysql
from mysql.connector.connection import MySQLConnection

from app import core_app
from app.core.settings import MariaDatabaseConfig


class MariaDatabaseWrapper:
    """
    MySQL Database
    """

    def __init__(self, db_config: MariaDatabaseConfig, handle_exception: bool = True):
        self.con = None
        self.error = None

        self.db_config = db_config
        self.handle_exception = handle_exception

    def __enter__(self):
        self.connect()
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        self.close()

    def close(self):
        if isinstance(self.con, MySQLConnection):
            self.con.close()

    def connect(self):
        """
        *   Establishes connection to the SQLite database file
        *   saves the connection instance in 'self.con'
        """
        try:
            self.con = mysql.connect(
                user=self.db_config.user,
                database=self.db_config.db_name,
                host=self.db_config.host,
                password=self.db_config.password
            )
        except Exception as ex:
            if self.handle_exception:
                core_app.logger.exception(ex)
                self.error = ex
            else:
                raise Exception(ex) from ex

    def get_cursor(self, prepared=False):
        """
        *   Can be called directly
        *   Helpful, when you want to 'Directly' interact with Sqlite3 module
        :param prepared:
        :return: Database Cursor object
        """
        self.connect()
        if not self.con:
            raise Exception("Unable to get db connection")
        return self.con.cursor(prepared=prepared)

    def modify_many(self, query, data):
        """

        :param query:
        :param data:
        :return:
        """
        updated_rows = 0
        try:
            if not data:
                return

            stmt = self.get_cursor(prepared=True)
            # Prepared Statement
            stmt.executemany(query, data)
            if not self.con:
                raise Exception("Unable to get db connection")
            self.con.commit()
            updated_rows = stmt.rowcount
        except Exception as ex:
            if self.handle_exception:
                core_app.logger.exception(ex)
                self.error = ex
            else:
                self.close()
                raise Exception(ex) from ex
        finally:
            self.close()
        return updated_rows

    def modify(self, query, *data):
        """
        *   Update or Insert rows in Main Database file
        *   Can also be used for creating Tables

        :param query: UPDATE SQL Query

        :return: Number of rows updated
        """
        updatedrows = 0
        try:
            if data:
                stmt = self.get_cursor(prepared=True)
                # Prepared Statement
                stmt.execute(query, data[0])
            else:
                stmt = self.get_cursor()
                stmt.execute(query)
            if not self.con:
                raise Exception("Unable to get db connection")

            self.con.commit()
            updatedrows = stmt.rowcount
        except Exception as ex:
            if self.handle_exception:
                core_app.logger.exception(ex)
                self.error = ex
            else:
                self.close()
                raise Exception(ex) from ex
        finally:
            self.close()
        return updatedrows

    def fetch(self, query, *data):
        """
        Run Sql Query and Fetch Rows

        :param query: SELECT SQL Query

        :return: Rows
        """

        rows = []
        try:
            if data:
                stmt = self.get_cursor(prepared=True)
                # Prepared Statement
                stmt.execute(query, data[0])
            else:
                stmt = self.get_cursor()
                stmt.execute(query)

            if stmt:
                rows = stmt.fetchall()
        except Exception as ex:
            if self.handle_exception:
                core_app.logger.exception(ex)
                self.error = ex
            else:
                self.close()
                raise Exception(ex) from ex
        finally:
            self.close()
        return rows

    def fetchone(self, query, *data):
        """
        Run Sql Query and fetches single row

        :param query: SELECT SQL Query

        :return: single row
        """
        row = None
        try:
            if data:
                # Prepared Statement
                stmt = self.get_cursor(prepared=True)
                stmt.execute(query, data[0])
            else:
                stmt = self.get_cursor()
                stmt.execute(query)

            if stmt:
                row = stmt.fetchone()
        except Exception as ex:
            if self.handle_exception:
                core_app.logger.exception(ex)
                self.error = ex
            else:
                self.close()
                raise Exception(ex) from ex
        finally:
            self.close()
        return row
