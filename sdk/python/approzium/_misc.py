import logging
import struct
from contextlib import contextmanager

logger = logging.getLogger(__name__)


def read_int32_from_bytes(bytes, index):
    num = struct.unpack("!i", bytes[index : index + 4])[0]
    return num


@contextmanager
def patch(module, attribute_name, new_value):
    original_value = getattr(module, attribute_name)
    setattr(module, attribute_name, new_value)
    try:
        yield
    finally:
        setattr(module, attribute_name, original_value)
