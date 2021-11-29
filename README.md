# pinterest-pin-creator

This projects helps me automate my pinterest pin creation. I can schedule the pins in the this [schudule file](./schedule.csv.example) and when running the application the schedule will be read and if a pin is up for creation, it will be created.

The project uses the [Pinterest API v5](https://developers.pinterest.com/docs/api/v5/) as well as the [quickstart guide](https://github.com/pinterest/api-quickstart#readme) as reference.

# Setup
To setup the project you first need prepare the `config.yaml` and `schedule.csv` files, as well as set the `APP_ID` and `APP_SECRET` environment variables.

## config.yaml setup
Execute the following commend to prepare your config.yaml file.
```
mv config.yaml.example config.yaml
```
Next fill in the configuration with your data:
```yaml
access_token_path: .access_token
schedule_file_path: /path/to/schedule.csv
browser_path: /path/to/a/browser/application
redirect_port: your_redirect_port
```
The redirect port must be the same that you set during your [Pinterest Application setup](https://developers.pinterest.com/docs/api/v5/#section/Configure-the-redirect-URI-required-by-this-code.)

## schedule.csv setup
Same as for your config file, you need to execute this command to prepare your schedule file.
```
mv schedule.csv.example schedule.csv
```

The schedule file has the following structure:

```csv
created;timestamp;board;title;description;filePath;link
```

- `created`: boolean flag - always fill in false
- `timestamp`: pin creation timestamp in the format 'Mon, 02 Jan 2006 15:04:05 MST'
- `board`: name of the pinterest board
- `title`: title for the pin
- `description`: description for the pin
- `filePath`: path to the image file for the pin
- `link`: link to the external URL of the pin

## App ID and App Secret
In order to create an access token you need to provide your [app ID and app secret](https://developers.pinterest.com/docs/api/v5/#section/Register-your-app-and-get-your-app-id-and-app-secret-key) as environment variables.

```
export APP_ID=<your app id>
export APP_SECRET=<your app secret>
```

# Running the code
While running the application you need to provide the `config.yaml` file as an argument.

```go
go run main.go config.yaml
```

# Building the application
```
go build -o ./bin/pin-creator  
```