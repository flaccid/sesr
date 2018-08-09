# sesr

Send email using Amazon SES by CLI or a simple web service.

## Usage

    $ sesr --help

### Environment Variables

These can be used as an alternative or in conjunction with the applicable CLI options.

- `AWS_REGION` - AWS region for the SNS service
- `AWS_ACCESS_KEY_ID` - AWS access key id
- `AWS_SECRET_ACCESS_KEY` - AWS secret access key

### Examples

#### CLI

```
sesr \
  --sender 'ACME Alerts <no-reply@acme.com>' \
  --recipients 'bill.gates@microsoft.com' \
  --subject 'you stole DOS' \
    'It took sometime to realize...'
```

#### Daemon Mode

Instead a CLI tool, run a RESTful API.

This needs to be run with the `--daemon` CLI option which exposes a http server.

```
POST /
{
  "sender": "ACME Alerts <no-reply@acme.com>",
  "recipients": [
    "bill.gates@microsoft.com"
  ],
  "subject": "you stole DOS",
  "body": "It took sometime to realize..."
}
```
