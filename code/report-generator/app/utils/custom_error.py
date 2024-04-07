

class AppError(Exception):
    """
    * Wrapper class for Exception to throw custom AppError

    Args:
        Exception (_type_): _description_
    """
    def __init__(self, message: str, status: int) -> None:
        self.message = message
        self.status = status
        super().__init__(self.message)