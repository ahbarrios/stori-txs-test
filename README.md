# Stori Challenge

In this repository you will find an example of **Stori** test exam to process a list of transaction and send an email with aggregated information extracted from it to a designated recipient using and *Outgoing* **SMTP** server. 

## Run

You can use **Docker** and `docker-compose` to run a local [test file](examples/txns.csv) or specify by your self. To run this application you will need the following set it environment variables:

- `STORI_SMTP_SERVER`: The outgoing **SMTP** server specified using **DNS** and port e.g `localhost:1025`
- `STORI_SMTP_USERNAME`: A *secret* username to log in into the **SMTP** server
- `STORI_SMTP_PASSWORD`: A *secret* password to log in into the **SMTP** server
- `STORI_SENDER`: The email address used to send this email with the **SMTP** server
- `STORI_RECIPIENT`: The recipient email address that will receive the resulting email with the summarized transaction information extracted from the provided file
- `STORI_SOURCE_PATH`: The full path in the system that stores the **CSV** transactions information

> NOTE: In this repo you will find a development **SMTP** [server](internal/email/server.go) that you can use for localstack.