# The Approzium CLI

## Usage

```
export APPROZIUM_SECRETS_MANAGER=local
export APPROZIUM_LOCAL_FILE_PATH=/path/to/secrets.yaml

approzium passwords list
```

```
export APPROZIUM_SECRETS_MANAGER=vault
export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=root

approzium passwords list
```

```
export APPROZIUM_SECRETS_MANAGER=asm
export AWS_REGION=us-east-1

approzium passwords list
```
