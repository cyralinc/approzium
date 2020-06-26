import boto3


def assume_role(iam_role):
    if iam_role is None:
        raise NotImplementedError("Automatic IAM ARN determination is not implemented")
    sts_client = boto3.client("sts")
    response = sts_client.assume_role(
        DurationSeconds=3600, RoleArn=iam_role, RoleSessionName="Service1",
    )
    return response


def obtain_credentials(response):
    return response["Credentials"]


def obtain_claimed_arn(response):
    return response["AssumedRoleUser"]["Arn"]


def obtain_signed_get_caller_identity(credentials):
    iam_session = boto3.Session(
        aws_access_key_id=credentials["AccessKeyId"],
        aws_secret_access_key=credentials["SecretAccessKey"],
        aws_session_token=credentials["SessionToken"],
    )
    iam_client = iam_session.client("sts")
    presigned_url = iam_client.generate_presigned_url("get_caller_identity", {})
    return presigned_url
