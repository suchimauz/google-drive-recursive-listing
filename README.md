# Listing files for Google Disk directory with uploading results to Google Spreadsheets with API
# Листинг файлов Google Диск в директории с загрузкой результата в Google таблицу через API

## Environments | Переменные окружения
```env
# You can take it from the URL, go to your Google drive folder and copy the id
# Можно взять с URL, заходите в свою папку на Google диске и копируете id
# https://drive.google.com/drive/folders/<someid>
ROOT_DIRECTORY_ID=<someid>

# You can take it from the URL, go to your table and copy the id
# Можно взять с URL, заходите в свою таблицу и копируете id
# https://docs.google.com/spreadsheets/d/<someid>/edit#gid=0
SPREADSHEET_ID=<someid>
```

## Prepare | Подготовка
Copy the .env.dist file and rename it to .env and edit the file, specifying ROOT_DIRECTORY_ID and SPREADSHEET_ID in it

Копируем файл .env.dist и переименовываем в .env и редактируем файл, указав в нем ROOT_DIRECTORY_ID и SPREADSHEET_ID

1
---
To get started, go to https://console.cloud.google.com/apis/ and create a new project. After creating the project, go to https://console.cloud.google.com/apis/library and enable the Google Drive API and Google Sheets API if they are not enabled.

Для начала, зайдите в https://console.cloud.google.com/apis/ и создайте новый проект. После создания проекта, переходим в раздел https://console.cloud.google.com/apis/library и включаем Google Drive API и Google Sheets API, если они не включены.

2
---
Next, create a Service Account, you can do this in the Credentials section: https://console.cloud.google.com/apis/credentials by clicking the + Create Credentials button.
We give any name and click Create and Continue, in the drop-down list "Select a Role" select Currently used -> Owner and click Continue, then click Done.

Далее следует создать Service Account, сделать это можно в разделе Credentials: https://console.cloud.google.com/apis/credentials, нажав кнопку + Create Credentials.
Даем любое название и нажимаем Create and Continue, в выпадающем списке "Select a Role" выбираем Currently used -> Owner и нажимаем Continue, после чего нажимаем Done.

3
---
Once you have created a Service Account, you need to create an OAuth 2.0 Client ID in https://console.cloud.google.com/apis/credentials, this can be done by clicking + Create Credentials, Application type, choose Desktop, any name. Next, you need to press the DOWNLOAD JSON button, rename the downloaded file to credentials.json and put it in the project folder.

Как только создали Service Account, нужно создать OAuth 2.0 Client ID в https://console.cloud.google.com/apis/credentials, это можно сделать, нажав + Create Credentials, Application type выбирайте Desktop, название любое. Далее нужно нажать кнопку DOWNLOAD JSON, скачанный файл переименовать в credentials.json и положить в папку проекта.

4
---
Now copy the E-Mail of the created Service Account. We go to the Google spreadsheet where the file listing will be uploaded. Click the "Access Settings" button and add the copied E-Mail, specifying the "Editor" role

Теперь копируем E-Mail созданного Service Account. Заходим в Google таблицу, куда будет выгружаться листинг файлов. Нажимаем кнопку "Настройки доступа" и добавляем скопированный E-Mail, указав роль "Редактор"

5
---
Now when I run the project with ```make run```, if you don't have a token configured yet or it's been a long time since the last run, you'll see the following message:

Теперь запускам проект командой ```make run```, если у вас еще не настроен токен или прошло много времени с момента последнего запуска, вы увидите следующее сообщение:

```
Go to the following link in your browser then type the
https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=<client_id>&redirect_uri=http%3A%2F%2Flocalhost&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fspreadsheets+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.metadata.readonly&state=state-token

authorization code:
```

We follow the link, select our account, click "Continue" and put all the checkmarks on the permission request and click "Continue" again, after which we will be redirected in the browser to a link of the form:

Переходим по ссылке, выбираем наш аккаунт, нажимаем "Продолжить" и ставим все галочки на запрос разрешения и еще раз нажимаем "Продолжить", после чего в браузере нас перенаправит на ссылку вида:

```
http://localhost/?state=state-token&code=<code>&scope=https://www.googleapis.com/auth/drive.metadata.readonly%20https://www.googleapis.com/auth/spreadsheets
```


We need to copy everything between code= and &, paste the resulting string ```<code>``` into the console and press Enter. After that the program will start, and if at the end you get the message: ```Done. Spreadsheet is successfully updated```, so your Google spreadsheet has been updated

Нам нужно скопировать все что находится между code= и &, вставить полученную строку ```<code>``` в консоль и нажать Enter. После чего программа запустится, и если в конце вы получили сообщение: ```Done. Spreadsheet is successfully updated```, значит ваша Google таблица обновилась

# Example | Пример
## At first run | При первом запуске
```
Go to the following link in your browser then type the
https://accounts.google.com/o/oauth2/auth?access_type=offline&client_id=<client_id>&redirect_uri=http%3A%2F%2Flocalhost&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fspreadsheets+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.metadata.readonly&state=state-token

authorization code: <code>
Saving credential file to: token.json
Имя                               Путь                               URL
test2                             /test2                             https://drive.google.com/drive/folders/...
тестовый файл.jpg                 /test2/тестовый файл.jpg           https://drive.google.com/file/d/...
test                              /test                              https://drive.google.com/drive/folders/...
1g64g62Rtnk.jpg                   /test/1g64g62Rtnk.jpg              https://drive.google.com/file/d/...
Выплаты по боту                   /Выплаты по боту                   https://docs.google.com/spreadsheets/d/...
6ZnCJfzPh-appointments.ndjson.gz  /6ZnCJfzPh-appointments.ndjson.gz  https://drive.google.com/file/d/...
8R6NxhND7-relatedperson.ndjson.gz /8R6NxhND7-relatedperson.ndjson.gz https://drive.google.com/file/d/...
65PnYq7q9-patients.ndjson.gz      /65PnYq7q9-patients.ndjson.gz      https://drive.google.com/file/d/...
6pfGgJhrv-concepts.ndjson.gz      /6pfGgJhrv-concepts.ndjson.gz      https://drive.google.com/file/d/...
62tZ8GqYV-codesystem.ndjson.gz    /62tZ8GqYV-codesystem.ndjson.gz    https://drive.google.com/file/d/...
6plKpdxLz-codesystem.ndjson.gz    /6plKpdxLz-codesystem.ndjson.gz    https://drive.google.com/file/d/...
6bVd4mlX8-lol.ndjson.gz           /6bVd4mlX8-lol.ndjson.gz           https://drive.google.com/file/d/...
8YNnQZpxG-codesystem.ndjson.gz    /8YNnQZpxG-codesystem.ndjson.gz    https://drive.google.com/file/d/...
Done. Spreadsheet is successfully updated
```

## At second run | При втором запуске:
```
Имя                               Путь                               URL
test2                             /test2                             https://drive.google.com/drive/folders/...
тестовый файл.jpg                 /test2/тестовый файл.jpg           https://drive.google.com/file/d/...
test                              /test                              https://drive.google.com/drive/folders/...
1g64g62Rtnk.jpg                   /test/1g64g62Rtnk.jpg              https://drive.google.com/file/d/...
Выплаты по боту                   /Выплаты по боту                   https://docs.google.com/spreadsheets/d/...
6ZnCJfzPh-appointments.ndjson.gz  /6ZnCJfzPh-appointments.ndjson.gz  https://drive.google.com/file/d/...
8R6NxhND7-relatedperson.ndjson.gz /8R6NxhND7-relatedperson.ndjson.gz https://drive.google.com/file/d/...
65PnYq7q9-patients.ndjson.gz      /65PnYq7q9-patients.ndjson.gz      https://drive.google.com/file/d/...
6pfGgJhrv-concepts.ndjson.gz      /6pfGgJhrv-concepts.ndjson.gz      https://drive.google.com/file/d/...
62tZ8GqYV-codesystem.ndjson.gz    /62tZ8GqYV-codesystem.ndjson.gz    https://drive.google.com/file/d/...
6plKpdxLz-codesystem.ndjson.gz    /6plKpdxLz-codesystem.ndjson.gz    https://drive.google.com/file/d/...
6bVd4mlX8-lol.ndjson.gz           /6bVd4mlX8-lol.ndjson.gz           https://drive.google.com/file/d/...
8YNnQZpxG-codesystem.ndjson.gz    /8YNnQZpxG-codesystem.ndjson.gz    https://drive.google.com/file/d/...
Done. Spreadsheet is successfully updated
```
