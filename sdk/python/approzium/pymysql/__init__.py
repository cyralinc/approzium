from ._connect import ApproziumConnection


def connect(*args, **kwargs):
    return ApproziumConnection(*args, **kwargs)


connect.__doc__ = ApproziumConnection.__init__.__doc__
__all__ = ["connect"]
