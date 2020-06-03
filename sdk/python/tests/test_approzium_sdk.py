import approzium
import approzium.misc
from approzium.misc import read_int32_from_bytes


def test_readint32():
    b = bytes([0, 0, 0, 0, 1, 0])
    assert read_int32_from_bytes(b, 0) == 0
    assert read_int32_from_bytes(b, 1) == 1
    assert read_int32_from_bytes(b, 2) == 256
