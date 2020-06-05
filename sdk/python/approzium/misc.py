import struct


def read_int32_from_bytes(bytes, index):
    num = struct.unpack("!i", bytes[index : index + 4])[0]
    return num
