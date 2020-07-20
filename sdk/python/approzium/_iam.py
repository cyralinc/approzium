from functools import lru_cache

import boto3
import requests
from ec2_metadata import ec2_metadata


def assume_role(iam_role):
    sts_client = boto3.client("sts")
    response = sts_client.assume_role(
        DurationSeconds=3600, RoleArn=iam_role, RoleSessionName="Service1",
    )
    return response


def obtain_credentials(response):
    return response["Credentials"]


def obtain_claimed_arn(response):
    return response["AssumedRoleUser"]["Arn"]


def obtain_signed_get_caller_identity(credentials=None):
    if credentials is None:
        # Attempt to use local identity through the boto3 client.
        sts_client = boto3.client("sts")
    else:
        sts_session = boto3.Session(
            aws_access_key_id=credentials["AccessKeyId"],
            aws_secret_access_key=credentials["SecretAccessKey"],
            aws_session_token=credentials["SessionToken"],
        )
        sts_client = sts_session.client("sts")
    presigned_url = sts_client.generate_presigned_url("get_caller_identity", {})
    return presigned_url


def get_local_arn():
    sts_client = boto3.client("sts")
    response = sts_client.get_caller_identity()
    return response["Arn"]


@lru_cache(maxsize=1)
def is_ec2():
    # AWS docs reference:
    # https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/identify_ec2_instances.html
    instance_identity_url = "http://169.254.169.254/latest/dynamic/instance-identity/"
    try:
        return requests.get(instance_identity_url, timeout=0.05).status_code == 200
    except requests.exceptions.ConnectionError:
        return False


def attribution_info():
    if not is_ec2():
        return {}
    attrs = ["public_hostname", "public_ipv4", "instance_id"]
    return {"ec2." + attr: getattr(ec2_metadata, attr) for attr in attrs}
