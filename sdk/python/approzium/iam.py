import boto3


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
