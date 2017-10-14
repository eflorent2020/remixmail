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

