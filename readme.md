# Muerta

Muerta - RESTful API for a term paper on "Web application to control the shelf life of products using computer vision".

## How to run?

First, create an `.env` file and put the following environment variables in it:

```shell
POSTGRES_USER=[psql_username]
POSTGRES_DB=[psql_db_name]
POSTGRES_PASSWORD=[psql_password]
API_PORT=[api_port]
SECRET=[api_secret]
RSA_PRIVATE_KEY=[api_rsa_private_key]
RSA_PUBLIC_KEY=[api_rsa_public_key]
```

Then Start the Docker containers with this command:

```shell
docker compose up -d --build
```

> Make sure you have open ports for the API and Database

## Features

- [ ] Python service to recognize shelf live on a picture
- [ ] JWT Authentication
- [ ] Logging in JSON format
- [ ] Users with roles and groups
