/*
Setting Up HTTPS for Your Go Application with Certbot and Let's Encrypt:

1. **Prerequisites**:
    - A registered domain name (for this guide, we're using `webhooks.example.com`).
    - A server with a public IP address where your application runs (and where you'll be running Certbot).
    - `certbot` tool installed. If not installed, you can get it from https://certbot.eff.org/.

2. **Point Your Domain to Your Server**:
    Make sure the domain `webhooks.example.com` points to the server's public IP address. Update the DNS settings of your domain to achieve this.

3. **Install Certbot**:
    The installation process varies depending on your server's OS. For Ubuntu, you can use:

    $ sudo apt-get update
    $ sudo apt-get install snapd
    $ sudo snap install --classic certbot
    $ sudo ln -s /snap/bin/certbot /usr/bin/certbot

4. **Obtain a Certificate**:

    a. **Using the Standalone Plugin**:
        If you don't have a web server running, you can use Certbot's standalone plugin to obtain a certificate:

        $ sudo certbot certonly --standalone

        Note: Certbot will bind to port 80, so ensure it's free.

5. **Locate the Certificate and Private Key**:
    After obtaining the certificate, Certbot will show the paths to your certificate and private key. Generally, they will be:

    - Certificate: `/etc/letsencrypt/live/webhooks.example.com/fullchain.pem`
    - Private Key: `/etc/letsencrypt/live/webhooks.example.com/privkey.pem`

    Use these paths in your Go application when calling `http.ListenAndServeTLS`.

6. **Automate Renewals**:
    Let's Encrypt certificates are valid for 90 days. It's recommended to renew them every 60 days. Certbot includes a cron job or systemd timer which will automatically renew the certificates.

    Test the renewal process with:

    $ sudo certbot renew --dry-run

7. **Integrate with Your Go Application**:

    Modify your Go server code to use:

    http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/webhooks.example.com/fullchain.pem", "/etc/letsencrypt/live/webhooks.example.com/privkey.pem", nil)

8. **Firewall Considerations**:
    Ensure ports 80 and 443 are open on your firewall to allow Certbot to validate your domain and for your application to serve over HTTPS.

9. **Important Notes**:
    - Keep your private key secret.
    - Back up `/etc/letsencrypt` regularly.
    - Always monitor the renewal process to ensure it's working as expected.

That's it! Your Go application should now be serving traffic over HTTPS using a Let's Encrypt certificate.

*/
