package collector

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

type channelsCollector []struct {
	val func(*channel) float64
	vec *prometheus.GaugeVec
}

// ChannelsCollector creates a new stats collector which is able to
// expose the channel metrics of a nsqd node to Prometheus. The
// channel metrics are reported per topic.
func ChannelsCollector(namespace string) StatsCollector {
	labels := []string{"type", "topic", "channel", "paused"}

	return channelsCollector{
		{
			val: func(c *channel) float64 { return float64(len(c.Clients)) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "client_count",
				Help:      "Number of clients",
			}, labels),
		},
		{
			val: func(c *channel) float64 { return float64(c.Depth) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "depth",
				Help:      "Queue depth",
			}, labels),
		},
		{
			val: func(c *channel) float64 { return float64(c.BackendDepth) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "backend_depth",
				Help:      "Queue backend depth",
			}, labels),
		},
		{
			val: func(c *channel) float64 { return float64(c.MessageCount) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "message_count",
				Help:      "Queue message count",
			}, labels),
		},
		{
			val: func(c *channel) float64 { return float64(c.InFlightCount) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "in_flight_count",
				Help:      "In flight count",
			}, labels),
		},
		{
			val: func(c *channel) float64 { return float64(c.DeferredCount) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "deferred_count",
				Help:      "Deferred count",
			}, labels),
		},
		{
			val: func(c *channel) float64 { return float64(c.RequeueCount) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "requeue_count",
				Help:      "Requeue Count",
			}, labels),
		},
		{
			val: func(c *channel) float64 { return float64(c.TimeoutCount) },
			vec: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "timeout_count",
				Help:      "Timeout count",
			}, labels),
		},
	}
}

func (coll channelsCollector) collect(s *stats, out chan<- prometheus.Metric) {
	for _, topic := range s.Topics {
		for _, channel := range topic.Channels {
			labels := prometheus.Labels{
				"type":    "channel",
				"topic":   topic.Name,
				"channel": channel.Name,
				"paused":  strconv.FormatBool(channel.Paused),
			}

			for _, c := range coll {
				c.vec.With(labels).Set(c.val(channel))
				c.vec.Collect(out)
			}
		}
	}
}
