import struct

MYSQLNativePassword = "mysql_native_password"


def get_auth_resp(
    authenticator, dbhost, dbport, dbuser, auth_plugin, auth_data, is_secure_connection
):
    plugin_auth_response = authenticator._get_mysql_hash(
        dbhost, dbport, dbuser, auth_plugin, auth_data
    )
    if is_secure_connection:
        resplen = len(plugin_auth_response)
        auth_response = struct.pack("<B", resplen) + plugin_auth_response
    else:
        auth_response = plugin_auth_response + b"\x00"
    return auth_response
