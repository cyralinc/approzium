---
title: Overview
---

import Image from '@theme/IdealImage';

## Introduction
Approzium enables developers to improve the observability and security of their applications. It allows applications to connect to databases without requiring access to credentials, and emits logs, metrics and traces with enriched information about their runtime execution context. It is built as a lightweight open source library with multi-language and multi-cloud support. 

Approzium eliminates blind spots in the diagnosis and tracing of complex performance problems within autoscaled microservices running on modern orchestration frameworks such as Kubernetes and AWS ECS (Elastic Container Service). For example, all instances of an autoscaled microservice look alike from a database’s perspective. This makes it harder to attribute performance issues resulting from misbehaving queries, buggy service code and faulty cloud VMs to a specific microservice instance.

By providing richer execution context about each microservice instance, such as the service’s IAM role info, the EC2 instance id and hostname where it’s running, or its container/task id in dockerized environments, Approzium allows DevOps teams to quickly trace and resolve performance issues. This context is added to the existing logs, metrics and traces already being emitted by the microservice.

Additionally, Approzium addresses common security vulnerabilities in how applications typically connect to a database. For example, here's how you might typically connect to your database currently:

```python
from psycopg2 import connect

conn = connect("host=1.2.3.4 user=user1 password=verySecurePassword")
```

Whether database credentials are stored in the application code itself (as in this example), or in a secrets manager such as Hashicorp Vault, allowing applications direct access to credentials exposes them to leaks through inadvertent application logging, application compromise, or theft of secrets manager API keys.

Approzium solves this problem by leveraging the cloud providers’ security infrastructure to authenticate the applications using IAM roles, instead, thus abstracting database credentials away from them. Here's what the same code would look like with Approzium 😎

```python
import approzium
from approzium.psycopg2 import connect

approzium.default_auth_client = approzium.AuthClient('authenticator:6001')
conn = connect("host=1.2.3.4 user=user1")  # Look ma, no password!
```

With the _password_ attribute removed from the `connect()` API, applications no longer need to know the actual database credentials in order to connect to them while administrators can still retain control over which applications are allowed access through the use of IAM roles.

Approzium integrates with popular secrets managers such as Hashicorp Vault and AWS Secrets Manager for the storage of database credentials, and uses OpenTelemetry to make it easy to integrate with a wide variety of observability products, such as Prometheus/Grafana, ELK, Datadog and New Relic.

Approzium was developed by the engineering team at Cyral and is available under the Apache 2.0 license, free for anyone to use and develop.

## Overview 
Approzium comprises two components that are deployed independently -- an SDK that runs as part of the application code, and a standalone service called the Authenticator with which the SDK interacts on behalf of the application. A single Authenticator instance can support multiple applications, which may be set up optionally as an auto scaling group for load balancing purposes.

<Image img={require('./images/overview-diagram.png')} />

Together, the SDK and Authenticator provide observability and security benefits to the application. The SDK has the ability to query infrastructure metadata services native to the various cloud platforms, and to orchestration frameworks such as Kubernetes to generate rich execution time context about the application. SDK APIs make it easy for developers to enrich their existing logs, metrics and traces with the context.

From a security perspective, the SDK also generates an application fingerprint that serves as its identity. While the actual mechanism depends on the cloud platform, the fingerprinting is primarily based on using an IAM role assigned to the application to generate a cryptographically verifiable signature for it. On AWS, for instance, Approzium SDK generates a signed GetCallerIdentity API request by requesting temporary credentials from the Security Token Service (STS). The signed request is sent to the Authenticator, which verifies its validity and authenticity by presenting it to STS.

The Authenticator service, in addition to authenticating applications, also computes responses to database challenges needed by the applications in order to connect to the databases. It is built with extensibility in mind to support multiple databases and their respective authentication methods. It delegates all credentials storage and management to popular secrets managers such as Hashicorp Vault and AWS Secrets Manager. By doing so, it ends up being a simple stateless service making it easy to handle failures and scale using auto scaling groups, Kubernetes daemon sets and ECS tasks.

## Next Steps
- You can get started with Approzium [here](quickstart).
- Learn how it works by checking out the [architecture](architecture).
- Visit our [GitHub](https://github.com/cyralinc/approzium) to download and play with it. PRs welcome!

We hope you enjoy it! 🤗

