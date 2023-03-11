# gsdb

gsdb stands for Google Sheet Database. It is a database service that is backed by Google Sheet.

(This repo is under active development. It is extremely beta. However, any contribution is much appreciated)

## Why is it needed?

The inspiration comes from hearing that [levels.fyi scaling to millions using Google Sheet](https://www.levels.fyi/blog/scaling-to-millions-with-google-sheets.html). There are many databases out there, most can run locally easily. However, provisioning the database in production often costs money and takes a while even if you are experienced. When you are running a startup or just want to quickly prototype an idea, you want a simple database that's quick, free, and easily queryable. In addition, business/analytics folks often prefer using Excel over SQL. Using gsdb saves you the headaches of exporting data from your database in a suitable format for the non-eng folks to digest.


### Google API Credential

(TODO: fill this out)

https://developers.google.com/drive/api/quickstart/go#set_up_your_environment
create service account, no role required

basically, follow https://www.labnol.org/google-api-service-account-220404
Then also go to https://console.cloud.google.com/apis/library/sheets.googleapis.com and click enable
