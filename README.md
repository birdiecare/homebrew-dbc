# `dbc`

**D**ata**b**ase **C**onnect

Connect to databases securely, utilising AWS SSM, EC2 & RDS IAM Authentication.

## Install `dbc`

`brew tap birdiecare/dbc`

`brew install dbc`

`dbc -h`

## Use `dbc`

The connect command opens a WebSocket port-forwarding session through an available Bastion instance.

Running `dbc connect` will start the process of creating a connection to a database.
Not specifying a host will prompt you to select one from a fuzzyfinder list of databases.

A connection to the specified Database host `-H` will then be available at `localhost:localport` (localport: `-lp`).

`dbc` picks up your configured AWS profile from your environment.

## Password Authentication

If the database you're connecting to doesn't have AWS IAM Authentication enabled, or doesn't have Users with the `rds_iam` role, you'll need to use a password to authenticate with the DB once your connection is open.

The following command will be enough to open a connection with a specified DB host with Birdie-specific defaults.

`dbc connect -H ${db_host}`

Once the connection is open:

```

➜ dbc connect -H ${db_host}
2023/03/21 17:27:16 DBConnect
2023/03/21 17:27:16 Using bastion: i-*
2023/03/21 17:27:16 Opening connection
2023/03/21 17:27:16 Connection Open at localhost:5432

...

```

You may connect to the database using `localhost:5432` as if `localhost` was the DB DNS.

`psql -h localhost -p 5432 -U ${user} -d ${db} --password`

## IAM Authentication

If your databases are super cool and secure, IAM Authentication will be enabled.

To use `dbc` with an IAM Enabled Database, you can use the `--iam` flag to generate a token while opening your connection!

```

➜ dbc connect -H ${db_host} --iam
2023/03/21 17:28:30 DBConnect IAM
2023/03/21 17:28:30 Token: ...
2023/03/21 17:28:30 Using bastion: i-*
2023/03/21 17:28:30 Opening connection
2023/03/21 17:28:30 Connection Open at localhost:5432

```

Then when connecting to your DB...

`psql -h localhost -p 5432 -U ${user} -d ${db} --password`

Paste the token!

Or... If you're very fancy:

`export PGPASSWORD=${token} && psql -h localhost -p 5432 -U ${user} -d ${db}`

## SOCKS5 proxy

`dbc` can also open a SOCKS5 proxy to allow any compatible application to connect to all databases through a proxy, thus not requiring to port-forward to a local port.

Simply login to the right environment, and run `dbc proxy`. `dbc` will open a SOCKS5 proxy on `localhost:1080` (port can be configured with `-p`).

Then, setup your database client's proxy settings to target the proxy port, leave `dbc proxy` running as long as you need, and enjoy!

When running Terraform plans locally for databases, you need to use this option. Otherwise, you would have to change your Postgres provider configuration (host and port) each time you need to plan locally. To use the SOCKS5 proxy with the Terraform PostgreSQL provider, set `ALL_PROXY=socks5://localhost:1080` before running your plan.

# Gotchas

## Ensuring your credentials

At times you may find that your use of DBC is interrupted by `sts caller-identity` errors. This is because `dbc` uses the `sts` endpoint to verify your AWS credentials before fetching any RDS endpoints or allowing any SSM sessions. If you’re finding yourself strugging with these, try the following steps:

Make sure you always log in with your sso profile e.g:

```bash
aws sso login --profile prod
```

Always export your profile, or use a tool like [granted](https://www.granted.dev/) to manage your profile setting:

```bash
export AWS_PROFILE=prod-read
```

Check your AWS caller identity outside of `dbc` first:

```bash
aws sts get-caller-identity
```

Ensure you have set your AWS region for your profiles:

```bash
vim ~/.aws/config
add a region to the necessary profiles
```

## Exit status `255` / `254`

If you're experiencing trouble opening a session, and you're recieving a `255` error, it's likely due to a missing AWS SSM Plugin installation.

Run this handy script! (Installs the plugin)

`wget -O - https://raw.githubusercontent.com/birdiecare/homebrew-dbc/main/install_ssm_plugin.sh | sh`

Or OSX:

`curl https://raw.githubusercontent.com/birdiecare/homebrew-dbc/main/install_ssm_plugin.sh | sh`
