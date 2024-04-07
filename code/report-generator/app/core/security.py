import re

def is_valid_table_name(table_name):
    if table_name and re.match(r"^[A-Za-z][a-zA-Z0-9._-]{1,40}$", table_name):
        return True
    return False
