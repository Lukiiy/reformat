package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type Meta struct {
    Pack *Pack `json:"pack"`
}

type Pack struct {
    Description any `json:"description"`
    PackFormat *int `json:"pack_format,omitempty"`
    SupportedFormats *[]int `json:"supported_formats,omitempty"`
    MinFormat *any `json:"min_format,omitempty"`
    MaxFormat *any `json:"max_format,omitempty"`
}

type Version struct {
    Major int
    Minor *int
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

    mainFormat, err := parse(os.Args[2])
    if err != nil {
        fmt.Println("Main format must be an integer")
        return
    }

    otherFormat := mainFormat

    if len(os.Args) >= 4 {
    	otherFormat, err = parse(os.Args[3])
        if err != nil {
            fmt.Println("Second format must be an integer")
            return
        }
    }

    min := mainFormat
    max := otherFormat

    if compare(min, max) > 0 {
        min, max = max, min
    }

    mode := detect(min.Major)

    meta.Pack.PackFormat = &mainFormat.Major
    meta.Pack.SupportedFormats = nil
    meta.Pack.MinFormat = nil
    meta.Pack.MaxFormat = nil

    switch mode {
        case "modern":
        	minVal := encode(min)
            maxVal := encode(max)

            meta.Pack.MinFormat = &minVal
            meta.Pack.MaxFormat = &maxVal
        case "transitional":
            meta.Pack.SupportedFormats = &[]int{min.Major, max.Major}
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

func parse(s string) (Version, error) {
    parts := strings.SplitN(s, ".", 2)

    major, err := strconv.Atoi(parts[0])
    if err != nil {
        return Version{}, err
    }

    if len(parts) == 1 {
        return Version{Major: major}, nil
    }

    minor, err := strconv.Atoi(parts[1])
    if err != nil {
        return Version{}, err
    }

    return Version{Major: major, Minor: &minor}, nil
}

func compare(first Version, second Version) int {
    if first.Major != second.Major {
        return first.Major - second.Major
    }

    firstMin := 0
    seconodMin := 0

    if first.Minor != nil {
        firstMin = *first.Minor
    }

    if second.Minor != nil {
        seconodMin = *second.Minor
    }

    return firstMin - seconodMin
}

func encode(version Version) any {
    if version.Minor == nil {
        return version.Major
    }

    return []int{version.Major, *version.Minor}
}