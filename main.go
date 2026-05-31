package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
)

type Meta struct {
    Pack *Pack `json:"pack"`
}

type Pack struct {
    Description any `json:"description"`
    PackFormat *int `json:"pack_format,omitempty"`
    SupportedFormats *[]int `json:"supported_formats,omitempty"`
    MinFormat *int `json:"min_format,omitempty"`
    MaxFormat *int `json:"max_format,omitempty"`
}

func main() {
    if len(os.Args) < 3 {
        fmt.Println("Missing arguments! Usage: <.mcmeta input> <main format> [secondary format]")
        return
    }

    file := os.Args[1]

    fileInfo, err := os.Stat(file)
    if err != nil {
        panic(err)
    }
    fileMode := fileInfo.Mode()

    stuff, err := os.ReadFile(file)
    if err != nil {
        panic(err)
    }

    var meta Meta
    if err := json.Unmarshal(stuff, &meta); err != nil {
        panic(err)
    }

    if meta.Pack == nil {
        return
    }

    mainFormat, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("Main format must be an integer")
        return
    }

    otherFormat := mainFormat

    if len(os.Args) >= 4 {
        otherFormat, err = strconv.Atoi(os.Args[3])
        if err != nil {
            fmt.Println("Second format must be an integer")
            return
        }
    }

    min := mainFormat
    max := otherFormat

    if min > max {
        min, max = max, min
    }

    mode := detect(min)

    meta.Pack.PackFormat = &mainFormat
    meta.Pack.SupportedFormats = nil
    meta.Pack.MinFormat = nil
    meta.Pack.MaxFormat = nil

    switch mode {
        case "modern":
            meta.Pack.MinFormat = &min
            meta.Pack.MaxFormat = &max
        case "transitional":
            meta.Pack.SupportedFormats = &[]int{min, max}
    }

    out, err := json.MarshalIndent(meta, "", "    ")
    if err != nil {
        panic(err)
    }

    err = os.WriteFile(file, out, fileMode)
    if err != nil {
        panic(err)
    }

    fmt.Println("Updated! Mode:", mode)
}

func detect(format int) string {
    if format > 64 { // 65 = 25w31a
        return "modern"
    }

    if format > 15 { // 16 = 23w31a
        return "transitional"
    }

    return "legacy"
}