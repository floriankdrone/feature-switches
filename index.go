package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "github.com/mattn/go-sqlite3"
)

type FeatureSwitch struct {
    ID     int    `json:"id"`
    Name   string `json:"name"`
    Value  bool   `json:"value"`  
}

var db *sql.DB

func main() {
    var err error
    db, err = sql.Open("sqlite3", "./feature_switches.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    createTable()

    http.HandleFunc("/", handleResources)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
    sqlStmt := `
        CREATE TABLE IF NOT EXISTS feature_switches (
            id INTEGER PRIMARY KEY,
            name TEXT,
            value INTEGER
        );
    `
    _, err := db.Exec(sqlStmt)
    if err != nil {
        log.Fatal(err)
    }
}

func handleFeatureSwitchs(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        getFeatureSwitchs(w, r)
    case "POST":
        createFeatureSwitch(w, r)
    case "PUT":
        updateFeatureSwitch(w, r)
    case "DELETE":
        deleteFeatureSwitch(w, r)
    default:
        http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
    }
}

func getFeatureSwitchs(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, value FROM feature_switches")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    resources := []FeatureSwitch{}
    for rows.Next() {
        var id int
        var name string
        err := rows.Scan(&id, &name)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        resources = append(resources, FeatureSwitch{ID: id, Name: name})
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resources)
}

func createResource(w http.ResponseWriter, r *http.Request) {
    var resource FeatureSwitch
    err := json.NewDecoder(r.Body).Decode(&resource)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := db.Exec("INSERT INTO resources(name) VALUES(?)", resource.Name)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    id, err := result.LastInsertId()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    resource.ID = int(id)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resource)
}

func updateResource(w http.ResponseWriter, r *http.Request) {
    var resource FeatureSwitch
    err := json.NewDecoder(r.Body).Decode(&resource)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := db.Exec("UPDATE resources SET name = ? WHERE id = ?", resource.Name, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "FeatureSwitch not found.", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func deleteResource(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(r.URL.Query().Get("id"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    result, err := db.Exec("DELETE FROM resources WHERE id = ?", id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if rowsAffected == 0 {
        http.Error(w, "FeatureSwitch not found.", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
