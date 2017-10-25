package influxdb

import (
	"github.com/influxdata/influxdb/client/v2"
)

type Batch struct {
	BP client.BatchPoints
}

type BatchConfig struct {
	Precision        string
	Database         string
	RetentionPolicy  string
	WriteConsistency string
}

func NewBatch(config BatchConfig) (Batch, error) {
	var bp client.BatchPoints
	var err error

	bp, err = client.NewBatchPoints(client.BatchPointsConfig{
		Precision:        config.Precision,
		Database:         config.Database,
		RetentionPolicy:  config.RetentionPolicy,
		WriteConsistency: config.WriteConsistency,
	})
	if err != nil {
		return Batch{}, err
	}

	batch := Batch{bp}

	return batch, nil
}

func (c Client) Write(batch Batch) error {
	return c.Conn.Write(batch.BP)
}
