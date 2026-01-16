package service

import (
	"bytes"
	"netboot/internal/config"
	"netboot/internal/model"
	"text/template"
)

type IPXEService struct {
	cfg *config.Config
}

func NewIPXEService(cfg *config.Config) *IPXEService {
	return &IPXEService{cfg: cfg}
}

const ipxeTemplate = `#!ipxe

# PXE Boot Menu
:start
menu PXE Boot Server - Select Image

{{range $i, $iso := .ISOs}}item iso{{$i}} {{$iso.Name}}
{{end}}
item shell iPXE Shell
item exit Exit
choose --default exit --timeout 30000 target && goto ${target}

{{range $i, $iso := .ISOs}}:iso{{$i}}
kernel http://${next-server}{{$.Port}}/isos/{{$iso.Path}}
boot

{{end}}
:shell
shell

:exit
exit
`

func (s *IPXEService) GenerateScript(isos []model.ISO) (string, error) {
	tmpl, err := template.New("ipxe").Parse(ipxeTemplate)
	if err != nil {
		return "", err
	}

	data := struct {
		ISOs []model.ISO
		Port string
	}{
		ISOs: isos,
		Port: s.cfg.Port,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
