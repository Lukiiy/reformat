package main

import (
    "encoding/json"
    "fmt"
    "os"
)

type Meta struct {
    Pack *Pack `json:"pack"`
}

type Pack struct {
    Description any `json:"description"`
    PackFormat *int `json:"pack_format,omitempty"`
    SupportedFormats *SupportedFormats `json:"supported_formats,omitempty"`
    MinFormat *int `json:"min_format,omitempty"`
    MaxFormat *int `json:"max_format,omitempty"`
}

type SupportedFormats struct {
    Min int `json:"min_inclusive"`
    Max int `json:"max_inclusive"`
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Missing arguments! Usage: <.mcmeta input> <main format> [compatible formats]")
        return
    }

    file := os.Args[1]

    stuff, err := os.ReadFile(file)
    if err != nil {
        panic(err)
    }

    var meta Meta
    if err := json.Unmarshal(stuff, &meta); err != nil {
        panic(err)
    }
}