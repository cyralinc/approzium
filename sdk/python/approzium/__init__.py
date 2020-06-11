import logging

authenticator_addr = "authenticator:1234"
iam_role = None

logging.basicConfig(
    format="[%(filename)s:%(lineno)s - %(funcName)10s() ] %(message)s",
)


def set_authenticator(new_authenticator):
    global authenticator_addr
    authenticator_addr = new_authenticator


def set_iam_role(new_iam_role):
    # this function tells Approzium to assume an IAM role and use that as its identity
    # with the authenticator service
    global iam_role
    iam_role = new_iam_role
