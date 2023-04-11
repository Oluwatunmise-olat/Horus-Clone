# ⚡️ Horus Clone

Clone of horus request logger created by `ichtrojon` with a twist 🔀. [horus](https://github.com/ichtrojan/horus)

### Features

-   Integrates with a `discord channel` which can trigger commands to get stats of endpoints logged to `horus-clone`.
-   Commands Extensions
-   Supports only mysql

### Example

```go

  import (
	  "fmt"
	  "log"
    "net/http"
	  db "github.com/Oluwatunmise-olat/Horus-Clone/db"
    clone "github.com/Oluwatunmise-olat/Horus-Clone"
  )

  func main(){
  // Mysql
  var conf db.Config = db.Config{ Password: "...", UserName: "...", Port: 3306, DatabaseName: "...", Host: "...", DiscordGuildId: "...", DiscordAppId: "...", DiscordToken: "..." }

	listener, err := clone.Init("mysql", &conf)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	// Serve Horus Clone
	if err := listener.Serve(":2024"); err != nil {
		fmt.Errorf(err.Error())
		return
	}

//  Register Routes
	http.Handle("/", listener.Watch(func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")

		response := map[string]any{"status": true, "message": "Clone is live 🧎", "data": nil}

		_ = json.NewEncoder(w).Encode(response)
	}))

	http.Handle("/api", listener.Watch(func (w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		response := map[string]any{"status": true, "message": "Clone api is live 🧎", "data": { "version": 1.0 }}

		_ = json.NewEncoder(w).Encode(response)
	}))


	log.Println("🚀 Server listening on port 8081")
	http.ListenAndServe(":8081", nil)
}

```

### TLDR

Just add the necessary discord and mysql envs, then access your registered routes (via browser/postman etc), then open your discord to the linked channel and try out `slash-commands` e.g `/welcome` or `/logs`
