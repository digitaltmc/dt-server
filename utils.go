package main

func CheckNilString(v *string) *string {
  if v == nil {
    ret := ""
    return &ret
  }
  return v
}

