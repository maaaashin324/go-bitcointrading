package models

import "time"

type DataFrameCandle struct {
	ProductCode string        `json:"product_code"`
	Duration    time.Duration `json:"duration"`
	Candles     []Candle      `json:"candles"`
}

func (d *DataFrameCandle) Times() []time.Time {
	times := make([]time.Time, len(d.Candles))
	for index, candle := range d.Candles {
		times[index] = candle.Time
	}
	return times
}

func (d *DataFrameCandle) Opens() []float64 {
	opens := make([]float64, len(d.Candles))
	for index, candle := range d.Candles {
		opens[index] = candle.Open
	}
	return opens
}

func (d *DataFrameCandle) Closes() []float64 {
	opens := make([]float64, len(d.Candles))
	for index, candle := range d.Candles {
		opens[index] = candle.Close
	}
	return opens
}

func (d *DataFrameCandle) Highs() []float64 {
	opens := make([]float64, len(d.Candles))
	for index, candle := range d.Candles {
		opens[index] = candle.High
	}
	return opens
}

func (d *DataFrameCandle) Lows() []float64 {
	opens := make([]float64, len(d.Candles))
	for index, candle := range d.Candles {
		opens[index] = candle.Low
	}
	return opens
}

func (d *DataFrameCandle) Volumes() []float64 {
	opens := make([]float64, len(d.Candles))
	for index, candle := range d.Candles {
		opens[index] = candle.Volume
	}
	return opens
}
