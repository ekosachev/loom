class HistoryException(BaseException):
    pass


class NoCurrentBranch(HistoryException):
    pass


class BranchDoesNotExist(HistoryException):
    pass


class BranchAlreadyExists(HistoryException):
    pass
