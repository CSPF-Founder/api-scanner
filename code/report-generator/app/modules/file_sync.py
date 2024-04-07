

import os
import shutil

from app.db.models.jobs import JobModel


def copy_from_remote(job: JobModel):
    """
    * Copy data from remote dir to local
    """
    src_dir = job.get_remote_work_dir()
    dest_dir = job.get_local_work_dir()

    if os.path.exists(dest_dir):
        shutil.rmtree(dest_dir)
    shutil.copytree(src_dir, dest_dir)
    return True


def copy_to_remote_only_log(job: JobModel):
    """
    * Copy only log to remote dir since scan failed
    """
    src_dir = job.get_local_work_dir()
    dest_dir = job.get_remote_work_dir()

    log_path = os.path.join(src_dir, "logs")
    dest_path = os.path.join(dest_dir, "logs")
    shutil.copy2(log_path, dest_path)


def copy_to_remote(job: JobModel):
    """
    * Copy reports,request archive and logs to remote dir
    """
    src_dir = job.get_local_work_dir()
    dest_dir = job.get_remote_work_dir()

    # copy report to remote dir
    report_path = os.path.join(src_dir, "report/report.docx")
    dest_report_path = os.path.join(dest_dir, "report.docx")
    shutil.copy2(report_path, dest_report_path)

    # copy request archive to remote dir
    request_path = os.path.join(src_dir, "requests")

    create_request_archive = False
    if os.path.exists(request_path) and os.listdir(request_path):
        create_request_archive = True
        
    if create_request_archive:
        dest_request_path = os.path.join(dest_dir, "error_requests")
        shutil.make_archive(dest_request_path, 'zip', request_path)

    # copy logs to remote dir
    log_path = os.path.join(src_dir, "logs")
    dest_log_path = os.path.join(dest_dir, "logs")
    shutil.copy2(log_path, dest_log_path)
