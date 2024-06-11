package main

import (
    "encoding/json"
    "io"
    "log"
	"fmt"
    "net/http"
    "os"
	"time"

    "path/filepath"
    "strconv"
    "strings"
    "github.com/fsnotify/fsnotify"
    "github.com/gorilla/websocket"
)

var (
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            // 允许来自任意来源的请求
            return true
        },
    }
    connections []*websocket.Conn
    repoPath    string
)

func main() {
    demoPath := os.Getenv("DEMO_PATH")
    if demoPath == "" {
        log.Fatal("DEMO_PATH environment variable is not set")
    }
    repoPath = filepath.Join(demoPath, "re2/repo")

    http.HandleFunc("/ws", handleWebSocket)
    http.HandleFunc("/queue", corsMiddleware(handlePost))
    http.HandleFunc("/seeds", corsMiddleware(getSeeds))
    http.HandleFunc("/seeds/", corsMiddleware(handleSeedRequest))

    // 使用 http.FileServer 和 http.StripPrefix 来处理静态文件请求
    fs := http.FileServer(http.Dir(repoPath))
    http.Handle("/", http.StripPrefix("/", fs))

    go watchFiles()

    port := 6767
    if len(os.Args) > 1 {
        var err error
        port, err = strconv.Atoi(os.Args[1])
        if err != nil {
            log.Fatalf("Invalid port number: %s", os.Args[1])
        }
    }

    server := &http.Server{Addr: ":" + strconv.Itoa(port), Handler: nil}
    log.Printf("Server started on http://localhost:%d", port)
    if err := server.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

func handlePost(w http.ResponseWriter, r *http.Request) {
    data, err := io.ReadAll(r.Body)
    if err != nil {
        log.Printf("Error reading request body: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    dirPath := filepath.Join(repoPath, "out/queue")
    if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
        log.Printf("Error creating directory: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fileName := getNextSeedFileName(dirPath)
    path := filepath.Join(dirPath, fileName)
    log.Printf("Handling POST request, saving file to: %s", path)
    if err = os.WriteFile(path, data, 0644); err != nil {
        log.Printf("Error saving file: %s", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    }
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    // 检查请求头
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    if r.Header.Get("Upgrade") != "websocket" {
        http.Error(w, "Upgrade header is missing or invalid", http.StatusBadRequest)
        return
    }
    if r.Header.Get("Connection") != "Upgrade" {
        http.Error(w, "Connection header is missing or invalid", http.StatusBadRequest)
        return
    }
    if r.Header.Get("Sec-Websocket-Version") != "13" {
        http.Error(w, "Unsupported WebSocket version", http.StatusBadRequest)
        return
    }

    // 升级为 WebSocket 连接
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("Error upgrading WebSocket connection: %s", err)
        return
    }
    defer func() {
        if err := conn.Close(); err != nil {
            log.Printf("Error closing WebSocket connection: %s", err)
        }
    }()

    connections = append(connections, conn)
    sendRepoContents(conn)

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading WebSocket message: %s", err)
            break
        }
    }
}

func watchFiles() {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        if err := watcher.Close(); err != nil {
            log.Printf("Error closing file watcher: %s", err)
        }
    }()

    done := make(chan bool)
    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    for _, conn := range connections {
                        sendRepoContents(conn)
                    }
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Printf("File watcher error: %s", err)
            }
        }
    }()

    log.Printf("Watching directory: %s", repoPath)
    if err = watcher.Add(repoPath); err != nil {
        log.Printf("Error watching directory: %s", err)
    }

    <-done
}

func sendRepoContents(conn *websocket.Conn) {
    var fileContents []map[string]string
    err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            relPath, err := filepath.Rel(repoPath, path)
            if err != nil {
                return err
            }
            content, err := os.ReadFile(path)
            if err != nil {
                return err
            }
            fileContents = append(fileContents, map[string]string{
                "name":    relPath,
                "content": string(content),
            })
        }
        return nil
    })
    if err != nil {
        log.Printf("Error walking directory: %s", err)
        return
    }

    log.Printf("Sending repo contents to connection %v", conn)

    // 检查 WebSocket 连接状态
    if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
        log.Printf("Error sending ping message: %s", err)
        return
    }

    // 发送文件内容
    if err := conn.WriteJSON(fileContents); err != nil {
        log.Printf("Error sending repo contents via WebSocket: %s", err)
        return
    }

    log.Printf("Finished sending repo contents to connection %v", conn)
}



func getSeeds(w http.ResponseWriter, r *http.Request) {
    log.Printf("Handling GET request for seeds with method: %s", r.Method)
    log.Printf("get seeds started")
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    seedDir := filepath.Join(repoPath, "in")
    log.Printf("Seed directory: %s", seedDir)
    files, err := os.ReadDir(seedDir)
    if err != nil {
        log.Printf("Error reading seed directory: %s", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    var seeds []map[string]interface{}
    for _, file := range files {
        if !file.IsDir() {
            content, err := os.ReadFile(filepath.Join(seedDir, file.Name()))
            if err != nil {
                log.Printf("Error reading seed file: %s", err)
                continue
            }
            seeds = append(seeds, map[string]interface{}{
                "name":    file.Name(),
                "content": string(content),
            })
            log.Printf("Seed file: %s, Content: %s", file.Name(), string(content))
        }
    }

    log.Printf("Seeds: %+v", seeds)

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(map[string]interface{}{
        "seeds": seeds,
    })
    if err != nil {
        log.Printf("Error encoding JSON response: %s", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
}

// func handleSeedRequest(w http.ResponseWriter, r *http.Request) {
//     body, err := io.ReadAll(r.Body)
//     if err != nil {
//         log.Printf("Error reading request body: %v", err)
//         http.Error(w, "Internal server error", http.StatusInternalServerError)
//         return
//     }
//     log.Printf("Request body: %s", string(body))

//     var requestData map[string]interface{}
//     if err := json.Unmarshal(body, &requestData); err != nil {
//         log.Printf("Error parsing request body: %v", err)
//         http.Error(w, "Invalid request payload", http.StatusBadRequest)
//         return
//     }

//     seedName, ok := requestData["name"].(string)
//     if !ok {
//         log.Printf("Invalid seed name in request body")
//         http.Error(w, "Invalid seed name", http.StatusBadRequest)
//         return
//     }

//     contentData, ok := requestData["content"].(map[string]interface{})
//     if !ok {
//         log.Printf("Invalid content in request body")
//         http.Error(w, "Invalid content", http.StatusBadRequest)
//         return
//     }

//     content, ok := contentData["value"].(string)
//     if !ok {
//         log.Printf("Invalid content value in request body")
//         http.Error(w, "Invalid content value", http.StatusBadRequest)
//         return
//     }

//     seedPath := filepath.Join(repoPath, "in", seedName)

//     switch r.Method {
//     case "PUT":
//         log.Printf("Updating seed: %s, Content: %s", seedName, content)
//         updateSeed(w, r, seedName, content)
//     case "DELETE":
//         deleteSeed(w, r, seedPath)
// 	case "POST":
//         handlePostRequest(w, r)	
//     default:
//         http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//     }
// }

func handleSeedRequest(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "PUT":
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Error reading request body: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        log.Printf("Request body: %s", string(body))

        var requestData map[string]interface{}
        if err := json.Unmarshal(body, &requestData); err != nil {
            log.Printf("Error parsing request body: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        seedName, ok := requestData["name"].(string)
        if !ok {
            log.Printf("Invalid seed name in request body")
            http.Error(w, "Invalid seed name", http.StatusBadRequest)
            return
        }

        contentData, ok := requestData["content"].(map[string]interface{})
        if !ok {
            log.Printf("Invalid content in request body")
            http.Error(w, "Invalid content", http.StatusBadRequest)
            return
        }

        content, ok := contentData["value"].(string)
        if !ok {
            log.Printf("Invalid content value in request body")
            http.Error(w, "Invalid content value", http.StatusBadRequest)
            return
        }

        log.Printf("Updating seed: %s, Content: %s", seedName, content)
        updateSeed(w, r, seedName, content)

    case "DELETE":
        seedName := filepath.Base(r.URL.Path)
        seedPath := filepath.Join(repoPath, "in", seedName)
        deleteSeed(w, r, seedPath)

    case "POST":
        body, err := io.ReadAll(r.Body)
        if err != nil {
            log.Printf("Error reading request body: %v", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
        }
        log.Printf("Request body: %s", string(body))

        var requestData map[string]interface{}
        if err := json.Unmarshal(body, &requestData); err != nil {
            log.Printf("Error parsing request body: %v", err)
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        content, ok := requestData["content"].(string)
        if !ok {
            log.Printf("Invalid content in request body")
            http.Error(w, "Invalid content", http.StatusBadRequest)
            return
        }

        log.Printf("Creating seed with content: %s", content)
        createSeed(w, r, content)

    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}




func updateSeed(w http.ResponseWriter, r *http.Request, seedName, content string) {
    seedPath := filepath.Join(repoPath, "in", seedName)
    
    if err := os.WriteFile(seedPath, []byte(content), 0644); err != nil {
        log.Printf("Error updating seed file: %s", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}

func deleteSeed(w http.ResponseWriter, r *http.Request, seedPath string) {
    if err := os.Remove(seedPath); err != nil {
        log.Printf("Error deleting seed file: %s", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
    var requestData map[string]string
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        log.Printf("Error parsing request body: %v", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    content, ok := requestData["content"]
    if !ok {
        log.Printf("Content not found in request body")
        http.Error(w, "Content is required", http.StatusBadRequest)
        return
    }

    createSeed(w, r, content)
}

func createSeed(w http.ResponseWriter, r *http.Request, content string) {
    seedDir := filepath.Join(repoPath, "in")
    if err := os.MkdirAll(seedDir, os.ModePerm); err != nil {
        log.Printf("Error creating seed directory: %s", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    seedName := "seed_" + time.Now().Format("20060102150405")
    seedPath := filepath.Join(seedDir, seedName)

    if err := os.WriteFile(seedPath, []byte(content), 0644); err != nil {
        log.Printf("Error creating seed file: %s", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    newSeed := map[string]string{
        "name":    seedName,
        "content": content,
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(newSeed)
}

func getNextSeedFileName(dirPath string) string {
    files, err := os.ReadDir(dirPath)
    if err != nil {
        log.Printf("Error reading directory: %s", err)
        return "seed"
    }

    var maxIndex int
    for _, file := range files {
        if !file.IsDir() && strings.HasPrefix(file.Name(), "seed") {
            indexStr := strings.TrimPrefix(file.Name(), "seed")
            if index, err := strconv.Atoi(indexStr); err == nil && index > maxIndex {
                maxIndex = index
            }
        }
    }

    return fmt.Sprintf("seed%d", maxIndex+1)
}