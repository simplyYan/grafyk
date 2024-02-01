package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
)

type Grafyk struct {
	elements map[string]interface{}
}

func New() *Grafyk {
	return &Grafyk{
		elements: make(map[string]interface{}),
	}
}

func (g *Grafyk) Progress(id string, value, max float64) {
	progressBar := NewProgressBar(value, max)
	g.elements[id] = progressBar
}

func (g *Grafyk) ProgressEdit(id string, value float64) {
	if progressBar, ok := g.elements[id].(*ProgressBar); ok {
		progressBar.SetValue(value)
	}
}

func (g *Grafyk) Graphic(id, filename string) error {
	graphic, err := NewGraphicFromTOML(filename)
	if err != nil {
		return err
	}
	g.elements[id] = graphic
	return nil
}

func (g *Grafyk) Echo(elements map[string]bool) {
	for id, element := range g.elements {
		if elements[id] {
			fmt.Println(element)
		}
	}
}

func (g *Grafyk) Destroy(id string) {
	delete(g.elements, id)
}

type ProgressBar struct {
	value, max float64
}

func NewProgressBar(value, max float64) *ProgressBar {
	return &ProgressBar{
		value: value,
		max:   max,
	}
}

func (p *ProgressBar) SetValue(value float64) {
	p.value = value
}

func (p *ProgressBar) String() string {
	progress := int((p.value / p.max) * 20)
	return fmt.Sprintf("[%s] %.0f%%", strings.Repeat("■", progress)+strings.Repeat("□", 20-progress), (p.value/p.max)*100)
}

type Graphic struct {
	data map[string]float64
	max  float64
}

func NewGraphic(data map[string]float64) *Graphic {
	max := 0.0
	for _, value := range data {
		if value > max {
			max = value
		}
	}
	return &Graphic{
		data: data,
		max:  max,
	}
}

func NewGraphicFromTOML(filename string) (*Graphic, error) {
	data := make(map[string]float64)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := toml.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return NewGraphic(data), nil
}

func (g *Graphic) String() string {
	var builder strings.Builder
	builder.WriteString("╭────────────────────────────────────────────────╮\n")

	for category, value := range g.data {
		progress := int((value / g.max) * 50)
		builder.WriteString(fmt.Sprintf("%3.0f%%  %s %s\n", (value/g.max)*100, strings.Repeat("█", progress), category))
	}

	builder.WriteString("╰────────────────────────────────────────────────╯")
	return builder.String()
}

func main() {
	gra := New()

	gra.Progress("myProgress", 0.453, 1.000)
	gra.ProgressEdit("myProgress", 0.800)

	if err := gra.Graphic("mainGrap", "graphicInfo.toml"); err != nil {
		fmt.Println("Error:", err)
	}

	elements := map[string]bool{
		"myProgress": true,
		"mainGrap":   true,
	}
	gra.Echo(elements)

	gra.Destroy("myProgress")
}
