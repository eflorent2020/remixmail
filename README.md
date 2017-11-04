# Purpose

This is an application for hidding email adresses of participants in a conversation.

The service give you an ephemaral alias so your recipient will not know your real adresses.

Exemple : like in Classified Ads service where seller and buyer 
communicate by mail together without disclosing their real address and thus trusting the service.

# Usage

This is a Go application for AppEngine Standard environement. 
Feel free to complete and adapt it to your needs. 
This software is copyright by myself https://github.com/emmanuel-florent licensed under the GPLV3.

[Setup a Go appengine project](https://cloud.google.com/appengine/docs/standard/go/download), edit app.yamp and config.Go

```bash
git clone https://github.com/emmanuel-florent/remixmail
cd remixmail 
npm run build
cd ..
gcloud app deploy app.yaml
```
# Mail API

visit the homepage for instructions on how to register a mail alias

# External API Usage

Get an API Key (api_key) for a user email (api_client) at /#/admin (admin login required)

then build header 
```bash
var hmac = crypto.createHmac('sha256', api_key);
hmac.update(request_path); // just the path, no urlencode
var signature = hmac.digest('hex'));
```
send request with header:
```bash
X-HMAC:  api_client,signature
```
Available endpoints are:
```bash
PUT /api/alias/{email}/{fullname}
DELETE /api/alias/{email}
```
# Devel
Back:
```bash
$ dev_appserver.py app.yaml
```
Front:
```bash
cd remixmail
remixmail$ npm run dev
```