variable "f5xc_ce_machine_image" {
  type = object({
    ingress_gateway = object({
        {{- range $key, $val := .ingress_gateway}}
        {{$key}} = string
        {{- end}}
    })
    ingress_egress_gateway = object({
        {{- range $key, $val := .ingress_egress_gateway}}
        {{$key}} = string
        {{- end}}
    })
  })
   default = {
      ingress_gateway = {
        {{- range $key, $val := .ingress_gateway}}
        {{$key}} = "{{$val.ami}}"
        {{- end}}
      }
      ingress_egress_gateway = {
        {{- range $key, $val := .ingress_egress_gateway}}
        {{$key}} = "{{$val.ami}}"
        {{- end}}
      }
    }
}