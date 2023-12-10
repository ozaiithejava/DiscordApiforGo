# DiscordApiforGo
a basic player - discord api json and html


## Creating Table in Mysql
```mysql
CREATE TABLE IF NOT EXISTS players (
    id INT AUTO_INCREMENT PRIMARY KEY,
    playername VARCHAR(255) NOT NULL,
    discord_id VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    coin INT NOT NULL,
    xp INT NOT NULL,
    level INT NOT NULL,
    eslenmedurumu BOOLEAN NOT NULL
);
```

## Dowland Go libs
```Bash
go get -u github.com/go-sql-driver/mysql
go get -u github.com/gorilla/mux
```

## Usage

 #### In Json Format:
   ```Link
http://localhost:3000/player?name=ur_name
```

 #### In Html Format:
```Link

http://localhost:3000/player?name=ur_name&inhtml=true ```