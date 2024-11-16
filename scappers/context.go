package scappers

import (
	"bytes"
	"html/template"

	"github.com/C-L-I-M/chapter-dong-dong/config"
	log "github.com/sirupsen/logrus"
)

var messageFormat = template.Must(template.New("message").Parse(`
{{.Name}} - {{.Number}}\n{{.Url}}
`))

type NewChapter struct {
	SagaSlug string
	Name     string
	Number   string
	Url      string
}

func (c NewChapter) String() string {
	var buf bytes.Buffer
	if err := messageFormat.Execute(&buf, c); err != nil {
		log.Error(err)
	}

	return buf.String()
}

type ScrappingContext struct {
	Saga    config.Saga
	state   map[string]any
	changed bool
}

func FromSaga(saga *config.Saga) *ScrappingContext {
	return &ScrappingContext{
		state: saga.State,
	}
}

func (c *ScrappingContext) GetState(key string) any {
	return c.state[key]
}

func (c *ScrappingContext) SetState(key string, value any) {
	c.state[key] = value
	c.changed = true
}

func (c *ScrappingContext) HasChanged() bool {
	return c.changed
}
