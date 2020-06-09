import boto3
import approzium


def obtain_signed_get_caller_identity():
    if approzium.iam_role is None:
        raise NotImplementedError('Automatic IAM ARN determination is not implemented yet')
    sts_client = boto3.client("sts")
    credentials = sts_client.assume_role(
        DurationSeconds=3600, RoleArn=approzium.iam_role, RoleSessionName="Service1",
    )["Credentials"]
    iam_session = boto3.Session(
        aws_access_key_id=credentials["AccessKeyId"],
        aws_secret_access_key=credentials["SecretAccessKey"],
        aws_session_token=credentials["SessionToken"],
    )
    iam_client = iam_session.client("sts")
    presigned_url = iam_client.generate_presigned_url("get_caller_identity", {})
    return presigned_url
