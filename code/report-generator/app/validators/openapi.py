from openapi_spec_validator import validate
from openapi_spec_validator.readers import read_from_filename
from urllib.parse import urlparse

# Check if openapi file is valid and has URL.


def validate_openapi_file(file_name_open_api):
    try:
        spec_dict, _ = read_from_filename(file_name_open_api)
        # Raises on an invalid spec.
        validate(spec_dict)
        url_from_spec = spec_dict['servers'][0]['url']
        if not url_from_spec:
            return False, None

        parse_result = urlparse(url_from_spec)
        if parse_result.scheme and parse_result.netloc:
            return True, url_from_spec
        return False, None
    except Exception as e:
        raise Exception("Validation Exception"+str(e))
