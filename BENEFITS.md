# I've tested for you: "the cloud without server"

Everybody is talking about the cloud, ...

but did you know that you could just put your app, without even having servers there ? 

This is called Serverless Cloud.

And *sometime* it can be a dream !

This type of service existed prior to cloud server hosting. My top three are :
- #3 AWS Lambda
- #2 Heroku
- #1 Google AppEngine

OK, let's say, i'm running an service that helps people meeets, and where they 
need to comunicate together by email, being able to reply without disclosing 
their real email, thus actually receiving an email, proxyied by the service.

Let's call this service RemixMail.

```
from: Emmanuel Florent <coof9Tia@servicenet>
to:	Jean-Pierre  <eeNgai1G@servicenet>
```
OK ?

OK !  let's envisage the three scenarios, self hosted, cloud hosted and serverless 
for our use case.

| Resource/Res.Cost | Self Hosted | Cloud Server | Serverless |
| ------------------| ----------- |------------- |----------- |
| Network           | very high   | low          | low        |
| Database Service  | moderate    | high         | low        |
| Web Service       | low         | high         | low        |
| Mail Service*     | very high*  | high*        | low        |
| Setup             | very high   | high         | moderate   |
| Sys. Administr.   | very high   | high         | low (none) |
| Sys.Maintenance   | very high   | moderate     | low (none) |
| Code              | very high   | high         | minimal    |
| Risk              | very high   | high         | minimal    |
| Human Resources   | very high   | moderate     | minimal    |

>* most often you'll outsource it in order to avoid a nightmare

The table show that cost of reliability remain high in standard cloud hosting for our use case.

I've designed a Serverless implementation in Go using Appengine: 

- Authentificated HTTP API
```
    - PUT /api/alias/:email/:fullname
    - GET /api/alias/:email
    - PUT /api/alias/validate/:key
    - DELETE /api/aliases/:email
```
- Incoming SMTP API
```
    -RCPT TO: anything@myservice
```

Receiving a registered email will :
- register the from adress after confirmation if needed
- transform the from/to email
- forward the email immediatly or after validation

# Conclusion

Serverless *may* provide huge benefits, but take care of :
- vendor lock-in
- language selection

# Readings

- [RFC 5322: Internet Message Format](https://www.ietf.org/rfc/rfc5322.txt)
- [https://cloud.google.com/appengine/docs/standard/go/datastore/](https://cloud.google.com/appengine/docs/standard/go/datastore/)
- [https://cloud.google.com/appengine/docs/standard/go/mail/sending-receiving-with-mail-api](https://cloud.google.com/appengine/docs/standard/go/mail/sending-receiving-with-mail-api)
- [https://www.reddit.com/r/serverless/](https://www.reddit.com/r/serverless/)
- Maciej Malawski.  "Towards Serverless Execution of Scientific Workflows HyperFlow Case Study",
In 11th Workshop on Workflows in Support of Large-Scale Science (WORKS@SC) , 
volume CEUR-WS 1800 of CEUR Workshop Proceedings , pages 25–33, Salt Lake City, Utah, USA, November
2016,[[http://ceur-ws.org/Vol-1800/paper4.pdf](http://ceur-ws.org/Vol-1800/paper4.pdf)]
- Jack Y. B. Lee and W. T. Leng, "Study of a Server-less Architecture for Video-on-Demand Applications",
In Proccedings of IEEE International Conference on Multimedia and Expo., Lausanne, Switzerland, 
26-29 August 2002, pp.233-236, [[http://www.mclab.info/lee_leung_icme2002.pdf](http://www.mclab.info/lee_leung_icme2002.pdf)]

