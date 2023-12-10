package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Player struct represents the player entity.
type Player struct {
	ID            int    `json:"id"`
	Playername    string `json:"playername"`
	DiscordID     string `json:"discord_id"`
	Status        string `json:"status"`
	Coin          int    `json:"coin"`
	XP            int    `json:"xp"`
	Level         int    `json:"level"`
	EslenmeDurumu bool   `json:"eslenmedurumu"`
}

var db *sql.DB

func main() {
	// MySQL bağlantısı açma
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/yourdatabase")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Router oluşturma
	router := mux.NewRouter()

	// API endpoint'leri tanımlama
	router.HandleFunc("/player", GetPlayer).Queries("name", "{name:[a-zA-Z0-9]+}").Methods("GET")
	router.HandleFunc("/player", GetPlayerHTML).Queries("name", "{name:[a-zA-Z0-9]+}", "inhtml", "{inhtml}").Methods("GET")

	// Web sunucu başlatma
	http.Handle("/", router)
	fmt.Println("Server is running on :3000")
	http.ListenAndServe(":3000", nil)
}

// GetPlayer, belirtilen oyuncu adı veya Discord ID'sine göre oyuncu bilgilerini döndürür.
func GetPlayer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]

	// MySQL sorgusu
	row := db.QueryRow("SELECT * FROM players WHERE playername = ? OR discord_id = ?", name, name)

	var player Player
	err := row.Scan(&player.ID, &player.Playername, &player.DiscordID, &player.Status, &player.Coin, &player.XP, &player.Level, &player.EslenmeDurumu)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	// JSON yanıtını oluşturma
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// GetPlayerHTML, belirtilen oyuncu adı veya Discord ID'sine göre oyuncu bilgilerini HTML olarak döndürür.
func GetPlayerHTML(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	inhtml := params["inhtml"]

	// MySQL sorgusu
	row := db.QueryRow("SELECT * FROM players WHERE playername = ? OR discord_id = ?", name, name)

	var player Player
	err := row.Scan(&player.ID, &player.Playername, &player.DiscordID, &player.Status, &player.Coin, &player.XP, &player.Level, &player.EslenmeDurumu)
	if err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	// HTML yanıtını oluşturma
	if strings.ToLower(inhtml) == "true" {
		renderHTML(w, player)
		return
	}

	// JSON yanıtını oluşturma
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// renderHTML, HTML ve CSS kullanarak oyuncu bilgilerini güzel bir şekilde gösterir.
func renderHTML(w http.ResponseWriter, player Player) {
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
	<head>
		<style>
			body {
				font-family: Arial, sans-serif;
				background-color: #f4f4f4;
				text-align: center;
			}
			.container {
				padding: 20px;
				background-color: #fff;
				border-radius: 10px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
				display: inline-block;
			}
			h2 {
				color: #333;
			}
			.player-info {
				text-align: left;
				margin-top: 20px;
			}
			.label {
				font-weight: bold;
				margin-right: 10px;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Player Information</h2>
			<div class="player-info">
				<p><span class="label">ID:</span>{{.ID}}</p>
				<p><span class="label">Playername:</span>{{.Playername}}</p>
				<p><span class="label">Discord ID:</span>{{.DiscordID}}</p>
				<p><span class="label">Status:</span>{{.Status}}</p>
				<p><span class="label">Coin:</span>{{.Coin}}</p>
				<p><span class="label">XP:</span>{{.XP}}</p>
				<p><span class="label">Level:</span>{{.Level}}</p>
				<p><span class="label">Eslenme Durumu:</span>{{.EslenmeDurumu}}</p>
			</div>
		</div>
	</body>
	</html>
	`

	tmpl, err := template.New("player").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Error rendering HTML", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, player)
}
