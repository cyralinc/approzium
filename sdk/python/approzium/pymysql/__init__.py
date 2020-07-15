from ._connect import ApproziumConnection


def connect(*args, **kwargs):
    return ApproziumConnection(*args, **kwargs)


__all__ = ['connect']
