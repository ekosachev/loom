class HistoryException(BaseException):
    pass


class NoCurrentBranch(HistoryException):
    pass


class BranchDoesNotExist(HistoryException):
    pass


class BranchAlreadyExists(HistoryException):
    pass


class WorkspaceException(BaseException):
    pass


class WorkspaceAlreadyExists(WorkspaceException):
    pass


class NoCurrentWorkspace(WorkspaceException):
    pass


class WorkspaceDoesNotExist(WorkspaceException):
    pass


class ConfigurationError(BaseException):
    pass


class NoActiveModel(ConfigurationError):
    pass


class ProviderException(BaseException):
    pass


class RateLimited(ProviderException):
    pass
