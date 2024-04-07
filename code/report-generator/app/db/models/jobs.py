
from datetime import datetime
from pydantic import BaseModel

import os

from app import core_app



class JobModel(BaseModel):
    id: int
    user_id: int
    status: int = 0
    created_at: datetime | None = None
    completed_time: datetime | None = None
    scanner_id: int | None = None
    scanner_ip: str | None = None

    def get_local_work_dir(self) -> str:
        return os.path.join(
            core_app.settings.local_temp_dir,
            f"user_{self.user_id}",
            f"job_{self.id}/"
        )

    def get_remote_work_dir(self) -> str:
        return os.path.join(
            core_app.settings.remote_work_dir,
            f"user_{self.user_id}",
            f"job_{self.id}/"
        )
