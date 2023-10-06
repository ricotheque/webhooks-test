# webhooks-test

Sample golang package that receives payloads from [Toggl Track API webhooks](https://developers.track.toggl.com/docs/webhooks_start). Includes payload validation.

## Requirements

For the package to run properly, the server must be running with HTTPS, with ports 80 and 443 open. On my development server, this was set up by using Certbot to create the certificate files, and `ufw` to open up the ports. 

A step-by-step on how to do this on an Ubuntu server is available [here](https://gist.github.com/ricotheque/5160387e6587f4a223369f131fab15fc).

## Setup overview
1. Enable HTTPS for the server
2. [Create the webhook subscription on Toggl Track](https://developers.track.toggl.com/docs/webhooks/subscriptions#post-creates-a-subscription)
3. Edit the config.yaml file to put in the certificate file locations and the webook subscription secret

## config.yaml

The package looks for `config.yaml` in the same directory. Here's an example `config.yaml`:

```yaml
certFile: /path/to/fullchain.pem
keyFile: /path/to/privkey.pem
togglWebhooks:
  secret: [secret]
```
