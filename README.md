## Qurandw

Download whole Quran chapters in one click

- Open the app
- Choose the reciter by inputting his number
- Wait till download is done

CLI options:

```bash

# create all directories/files in English
# defaults to ar - [ISO 639-1 language code]
./qurandw en 

# specify a download directory
# note that a subdirectory will be created for the reciter
# defaults to current work directory
./qurandw en ./downloads 

# download small chapters first (Quran in reverse order)
# defaults to false
./qurandw en ./downloads true

# download in batches of 20, defaults to 10
./qurandw en ./downloads true 20

```

Note that on *NIX systems, you'll need to set proper file permissions

```shell
chmod +x ./qurandw-darwin-arm64
```

Then open the app

```shell
./qurandw-darwin-arm64
```