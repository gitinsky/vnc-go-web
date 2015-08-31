package main

func ternary(c bool, t interface{}, f interface{}) interface{} {
    if c {
        return t
    }
    return f
}