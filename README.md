# trello-burndown
An easy to use trello burndown chart.

### Screenshots

#### Viewing a burndown chart
![view](screenshots/view.png)

#### Adding a trello board
![add](screenshots/add.png)

#### Index page with table of trello boards
![index](screenshots/index.png)

### Assumptions

- The points must be present in the title between parenthesis like so: `(2) Add login page`
- The last column of the board is where finished cards are found.

### Installation

#### Obtain trello tokens
1. Login to [trello](https://trello.com)
2. [Generate a Developer API key](https://trello.com/app-key)
3. Generate a token by visiting the following URL:
`https://trello.com/1/authorize?name=trello-burndown&expiration=never&response_type=token&key=DEVELOPER_API_KEY`.
Replace "DEVELOPER_API_KEY" with the key you generated in the previous step.
4. Write both the Developer API key and the generated token down, you will need these to configure the application.

#### Docker: Setup & Running
1. Create a new directory to store the configuration and Sqlite3 database.

    ```
    λ mkdir trello-burndown && cd trello-burndown
    ```

2. Create a file named `config.yaml` in the same directory, copy the contents from the default [here](https://github.com/swordbeta/trello-burndown/blob/master/config.yaml.default).
3. Edit the configuration file with your favorite editor and set the developer api key and generated token you wrote down earlier.
4. Run it! (The config file must be present in the /root directory inside the docker container.)

    ```
    λ docker run -d -p 8080:8080 --volume $(pwd):/root/ swordbeta/trello-burndown:v1.0.0
    ```

#### Binary: Setup & Running
1. Download the latest release from [here](https://github.com/swordbeta/trello-burndown/releases).
2. Create a file named `config.yaml` in the same directory, copy the contents from the default [here](https://github.com/swordbeta/trello-burndown/blob/master/config.yaml.default).
3. Edit the configuration file with your favorite editor and set the developer api key and generated token you wrote down earlier.
4. Run it! You could run this as a daemon with upstart/supervisord/systemctl/etc.

    ```
    λ ./trello-burndown
    ```
    
