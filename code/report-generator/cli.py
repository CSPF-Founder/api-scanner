"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""
import argparse

from app.controllers.main_controller import MainController


def cli():
    """
    * Main CLI function
    """
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-m",
        "--module",
        help="Module to run",
        dest="module",
        required=True,
        choices=["reporter"],
    )
    parser.add_argument("--jobid", "-j", type=str, required=True, dest="job_id")
    parser.add_argument("--userid", "-u", type=str, required=True, dest="user_id")

    args = parser.parse_args()

    if not args:
        parser.print_help()
        return

    try:
        if args.module == "reporter":
            controller = MainController(args.job_id, args.user_id)
            controller.run()
        else:
            parser.print_help()
    except Exception as ex:  # pylint: disable=broad-except
        print(ex)


if __name__ == "__main__":
    cli()
