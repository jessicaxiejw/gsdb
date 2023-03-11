# gsdb

gsdb stands for Google Sheet Database. It is a database service that is backed and stored using Google Sheet.

(This repo is under active development. It is extremely beta. However, any contribution is much appreciated)

## Why is it needed?

The inspiration comes from hearing that [levels.fyi scaling to millions using Google Sheet](https://www.levels.fyi/blog/scaling-to-millions-with-google-sheets.html). There are many databases out there, most can run locally easily. However, provisioning the database in production often costs money and takes a while even if you are experienced.

When you are running a startup or just want to quickly prototype an idea, you want a simple database that's quick, free, and easily queryable. In addition, business/analytics folks often prefer using Excel over SQL. Using gsdb saves you the headaches of exporting data from your database in a suitable format for the non-eng folks to digest.


## Setup

In order to use gsdb, you need to first create a folder in your own Google Drive and then share that with a service account created in the Google Cloud Console. If you are struggling with any of the steps below, check out Amit's guide on [How to Upload Files to Google Drive with a Service Account](https://www.labnol.org/google-api-service-account-220404).

### <a name="createfolder"></a>Create a folder for sharing in Google Drive

1. Go to [Google Drive](https://drive.google.com/drive/u/0/my-drive) and make sure you have the right Google account selected.
1. Select "My Drive" on the left hand side.
1. Click "New" right above "My Drive" and select "New folder". You can name the folder anything (eg/ you can name it `gsdb`). This will be the root directory which your new database will be stored in.

### <a name="createserviceaccount">Create a service account

Now we will create a Google Cloud service account and enable the necessary APIs:

1. Go to [Google Drive API](https://console.cloud.google.com/apis/api/drive.googleapis.com) and click `ENABLE`. You may need to go through the set up of a new Google Cloud Platform project if your organization doesn't have one.
1. Go to [Google Sheet API](https://console.cloud.google.com/apis/library/sheets.googleapis.com) and click `ENABLE`.
1. Now that both APIs have been enabled, let's create the service account by going to [the credentials page](https://console.cloud.google.com/apis/credentials) under Google Cloud API & Services.
1. Click `+ CREATE CREDENTIALS` near the top of the credentials page.
1. Choose `Service account` when given different types of credentials.
1. Fill out the `Service account details` and leave the rest blanket. This step will only affects the name of the service account. Then click `DONE`.
1. Now we can generate a new credential JSON file by going to [Google Cloud API & Services -> Credentials](https://console.cloud.google.com/apis/credentials). Under `Service Accounts`, you should see the newly created account in the previous step.
1. Click on the service account you just generated, then `KEYS` -> `ADD KEY` -> `Create new key` -> `JSON`. Complete the key creation by clicking `CREATE`. Save the newly generated JSON key in a file (eg/ in `credential.json`) and move it to a secure place. You will need it later.
1. In order for the newly created service account to access the folder created by you, you must share the permission. Right click on the folder created in Google Drive in [the previous section](#createfolder), choose `Share`, and share it with the email address of the service account. The permission should be `editor`.

## Run

If you haven't finished the setup section above, please complete it first.

To run the application, do
```
git clone https://github.com/jessicaxiejw/gsdb
cd gsdb
go run main.go server -c <path to the JSON credential file of your service account> // eg/ go run main.go server -c credential.json
```

I am planning on releasing the binary in brew and other linux distros. However, for now, you will have to run the go binary.
