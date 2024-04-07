from app.db.models.jobs import JobModel

from app.utils import common_utils


class CrudJob:
    TABLE_NAME = "jobs"

    def __init__(self, db_session):
        self.db_session = db_session

    def find_by_id(self, job_id: int) -> JobModel | None:
        """
        Find entry by id
        """
        query = (
            f"SELECT id, user_id, status, created_at FROM {self.TABLE_NAME} "
            f"WHERE id = %s"
        )
        data = (job_id,)

        row = self.db_session.fetchone(query, data)
        if not row:
            return None

        id = common_utils.convert_to_int(row[0])
        if not id:
            return None
        user_id = common_utils.convert_to_int(row[1])
        if not user_id:
            return None
        status = common_utils.convert_to_int(row[2])
        if not status:
            return None
        created_at = common_utils.convert_to_datetime(row[3])

        return JobModel(id=id, user_id=user_id, status=status, created_at=created_at)
