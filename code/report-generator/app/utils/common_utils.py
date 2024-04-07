"""
Copyright (c) 2022 CySecurity Pte. Ltd. - All Rights Reserved
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by CySecurity Pte. Ltd.
"""
from datetime import datetime
import re


def smart_str(input_text):
    """
    * Proper handling for string
    * Converts any input like unicode,numbers into string value

    Args:
        input_text (_type_): _description_

    Returns:
        _type_: _description_
    """
    if not input_text:
        return input_text

    if isinstance(input_text, str):
        return input_text
    if isinstance(input_text, (bytearray, bytes)):
        return str(input_text, "utf-8")
    if isinstance(input_text, (int, float)):
        return str(input_text)

    return str(input_text, 'utf-8')


def convert_to_int(input_text) -> int:
    """
    *   Converts any input like unicode,numbers into int value
    *   If conversion fails, raises ValueError

    Args:
        input_text (_type_): _description_

    Returns:
        _type_: _description_
    """
    if input_text is None:
        raise ValueError("Input cannot be None")

    if isinstance(input_text, int):
        return input_text

    if isinstance(input_text, str):
        try:
            return int(input_text)
        except ValueError:
            raise ValueError("Input is not a valid integer")

    raise ValueError("Input is not a valid integer")


def convert_to_int_or_none(input_text):
    """
    *   Converts any input like unicode,numbers into int value
    *   If conversion fails, returns None

    Args:
        input_text (_type_): _description_

    Returns:
        _type_: _description_
    """
    if input_text is None:
        return None

    if isinstance(input_text, int):
        return input_text

    if isinstance(input_text, str):
        try:
            return int(input_text)
        except ValueError:
            return None

    return None


def convert_to_datetime(input_date) -> datetime:
    """
    *   Converts any input like unicode,numbers into datetime value
    *   If conversion fails, raises ValueError

    Args:
        input_date (_type_): _description_

    Returns:
        _type_: _description_
    """
    if not input_date:
        raise ValueError("Input cannot be None")

    if isinstance(input_date, datetime):
        return input_date

    if isinstance(input_date, str):
        try:
            return datetime.strptime(input_date, "%Y-%m-%d %H:%M:%S")
        except ValueError:
            raise ValueError("Input is not a valid datetime")

    raise ValueError("Input is not a valid datetime")


def replace_unsupported_for_docx(input_value: str, replace_with="[REFER_REQUEST_FILE]"):
    # Function to remove unsupported characters from a string
    # Match any character outside the BMP (Basic Multilingual Plane)
    pattern = re.compile(r'[^\u0000-\uffff]')
    return pattern.sub(replace_with, input_value)


def any_unsupported_for_docx(input_value: str) -> bool:
    # Match any character outside the BMP (Basic Multilingual Plane)
    pattern = re.compile(r'[^\u0000-\uffff]')
    return pattern.search(input_value) is not None
