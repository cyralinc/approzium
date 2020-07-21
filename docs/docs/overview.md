---
title: Overview
---

import Image from '@theme/IdealImage';

## Introduction
Approzium enables developers to improve the observability and security of their applications. It allows applications to connect to databases without requiring access to credentials, and emits logs, metrics and traces with enriched information about their runtime execution context. It is built as a lightweight open source library with multi-language and multi-cloud support.

Approzium eliminates blind spots in the diagnosis and tracing of complex performance problems within autoscaled microservices running on modern orchestration frameworks such as Kubernetes and AWS ECS (Elastic Container Service). It makes it easy to attribute performance issues resulting from misbehaving queries, buggy service code and faulty cloud VMs to a specific microservice instance.

Additionally, Approzium addresses common security vulnerabilities in how applications typically connect to a database. For example, here's how you might typically connect to your database currently:

```python
from psycopg2 import connect

# In your environment, the password may be pulled from any place, S3, a config 
# file, HashiCorp Vault, or wherever. But wherever it originates, it can be
# leaked by your client application because it's in memory.
conn = connect('host=1.2.3.4 user=user1 password=verySecurePassword')
```

Whether database credentials are stored in the application code itself (as in this example), or in a secrets manager such as Hashicorp Vault, allowing applications direct access to credentials exposes them to leaks through inadvertent application logging, application compromise, or theft of secrets manager API keys.

Approzium solves this problem by leveraging the cloud providers’ security infrastructure to authenticate the applications using IAM roles. Here's what the same code would look like with Approzium 😎

```python
import approzium
from approzium.psycopg2 import connect

approzium.default_auth_client = approzium.AuthClient('authenticator:6001')
conn = connect('host=1.2.3.4 user=user1')  # Look ma, no password!
```

With the _password_ attribute removed from the `connect()` API, applications no longer need to know the actual database credentials in order to connect to them while administrators can still retain control over which applications are allowed access through the use of IAM roles.

Approzium is developed by the engineering team at Cyral and is available under the Apache 2.0 license, free for anyone to use and develop.


## Overview 
Approzium comprises two components that are deployed independently -- an SDK that runs as part of the application code, and a standalone service called the Authenticator with which the SDK interacts on behalf of the application.

<Image img={require('./images/overview-diagram.png')} />

Together, the SDK and Authenticator provide observability and security benefits to the application. The SDK has the ability to query infrastructure metadata services native to the various cloud platforms, and to orchestration frameworks such as Kubernetes to generate rich execution time context about the application.

The SDK also generates an application fingerprint which serves as its identity. Fingerprinting is primarily based on using an IAM role assigned to the application to generate a cryptographically verifiable signature. On AWS, Approzium SDK generates a signed _GetCallerIdentity_ API request using the _Security Token Service (STS)_, which is sent to the Authenticator for verification.

The Authenticator also computes responses to database challenges on behalf of applications. All credentials storage and management is delegated to secrets managers such as Hashicorp Vault. By doing so, it ends up being a simple stateless service making it easy to handle failures and scale using auto scaling groups, Kubernetes daemon sets and ECS tasks.

## Next Steps
- You can get started with Approzium [here](quickstart).
- Learn how Approzium works by checking out the [architecture](architecture).
- Learn more about Approzium's approach to [observability](observability) and [security](security-model).
- Visit our [GitHub](https://github.com/cyralinc/approzium) to download and play with it. PRs welcome!
- Join our [Slack](https://join.slack.com/t/approzium/shared_invite/zt-fg9bdcfa-H9YFnlg3XeosKyMIYadmcg) for help and announcements.

We hope you enjoy it! 🤗

