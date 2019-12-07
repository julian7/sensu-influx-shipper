package tcpserver

import (
	"errors"
	"io"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/influxdata/influxdb/client/v2"
)

type Conn struct {
	log.Logger
	client.Client
	Grouping bool
	DB       string
}

func (c *Conn) handle(reader io.Reader) {
	defer func() { _ = c.Log("msg", "connection closed") }()

	_ = c.Log("msg", "new connection")

	event, err := ReadEvent(reader)
	if err != nil {
		c.logError(err)
		return
	}

	bp, err := c.BatchMetrics(event)
	if err != nil {
		c.logError(err)
		return
	}

	if err := c.Client.Write(bp); err != nil {
		c.logError(NewError("write metric point", err))
	}
}

func (c *Conn) BatchMetrics(event *Event) (client.BatchPoints, error) {
	baseTags := event.BaseTags()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  c.DB,
		Precision: "s",
	})
	if err != nil {
		return nil, NewError("set batch points", err)
	}

	for _, point := range event.Metrics.Points {
		var key string

		tags := map[string]string{}

		for key, val := range baseTags {
			tags[key] = val
		}

		for _, tag := range point.Tags {
			tags[tag.Name] = tag.Value
		}

		name, key, metric := splitName(point.Name, c.Grouping)

		if len(metric) > 0 {
			tags["metric"] = metric
		}

		fields := map[string]interface{}{key: point.Value}

		pt, err := client.NewPoint(name, tags, fields, time.Unix(point.Timestamp, 0))
		if err != nil {
			c.logError(NewErrorWithValue("create metric point", point.Name, err))
		}

		bp.AddPoint(pt)
	}

	return bp, nil
}

func (c *Conn) logError(err error) {
	data := []interface{}{}
	actionErr := &Error{}

	if errors.As(err, &actionErr) {
		data = append(data, "action", actionErr.Action())
		if actionErr.Value() != nil {
			data = append(data, "value", actionErr.Value())
		}
	}

	data = append(data, "error", err)
	_ = c.Log(data...)
}
