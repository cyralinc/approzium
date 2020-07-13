# The Approzium Authenticator

The Approzium authenticator is a server that stands between a client, and a database it wants access to. It verifies
the client's identity using the _platform_ in which it is running, like by using its AWS identity credentials. Then, it 
checks a place where database credentials are stored, like in HashiCorp Vault, builds a connection, and passes the 
_connection_ back to the client. So, the client never sees the database password. And thus, the client _can't leak_
the password.

## Developing Is Easy!

We love contributions. To easily develop, in the `authenticator` folder, run `$ make dev`. Then, run the authenticator.

```
$ authenticator -dev
```

It will start the authenticator up on your `localhost` without TLS. Check that it's up by hitting its API.

```
$ curl -v http://localhost:6000/v1/health
```

For test data, simply edit the `secrets.yaml` file the authenticator's logs mention, and restart it. Or, use the test
data that's already there.
