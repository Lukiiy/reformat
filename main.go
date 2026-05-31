package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "slices"
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
        fmt.Println("Missing arguments! Usage: <.mcmeta input> <main format> [compatible formats]")
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

    allFormats := []int{mainFormat}

    for _, arg := range os.Args[3:] {
        val, err := strconv.Atoi(arg)
        if err != nil {
            fmt.Println("Compatible formats must be all integers")
            return
        }

        allFormats = append(allFormats, val)
    }

    mode := detect(meta)

    meta.Pack.PackFormat = &mainFormat
    meta.Pack.SupportedFormats = nil
    meta.Pack.MinFormat = nil
    meta.Pack.MaxFormat = nil

    min := slices.Min(allFormats)
    max := slices.Max(allFormats)

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

func detect(meta Meta) string {
    if meta.Pack.MinFormat != nil || meta.Pack.MaxFormat != nil {
        return "modern"
    }

    if meta.Pack.SupportedFormats != nil {
        return "transitional"
    }

    return "legacy"
}