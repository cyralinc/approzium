from approzium._postgres import construct_msg, parse_msg


def test_parse_msg():
    # empty messages with a header
    msg = b"R\x00\00\00\04"
    assert parse_msg(msg) == (ord("R"), b"")
    msg = b"Z\x00\00\00\04"
    assert parse_msg(msg) == (ord("Z"), b"")
    # message with one byte content
    msg = b"M\x00\00\00\05\13"
    assert parse_msg(msg) == (ord("M"), b"\13")
    # message with one byte content but lots of extra content
    # extra content should be ignored
    msg = b"M\x00\00\00\05\13" + 100 * b"E"
    assert parse_msg(msg) == (ord("M"), b"\13")


def test_construct_msg():
    assert construct_msg(b"R", b"") == b"R\x00\00\00\04"
    assert construct_msg(b"Z", b"") == b"Z\x00\00\00\04"
    assert construct_msg(b"R", b"\x0c") == b"R\x00\00\00\05\x0c"


def test_integration_parse_construct_msg():
    headers = [b"R", b"Z"]
    contents = [b"testbytestring", b""]
    for header, content in zip(headers, contents):
        msg = construct_msg(header, content)
        assert parse_msg(msg) == (ord(header), content)
