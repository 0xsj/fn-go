// services/user-service/pkg/metrics/metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// UserCreationCounter counts user creations
	UserCreationCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_creation_total",
		Help: "The total number of users created",
	})

	// UserCreationErrorCounter counts errors during user creation
	UserCreationErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_creation_errors_total",
		Help: "The total number of errors during user creation",
	})

	// UserFetchCounter counts user fetch operations
	UserFetchCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_fetch_total",
		Help: "The total number of user fetch operations",
	})

	// UserFetchErrorCounter counts errors during user fetch operations
	UserFetchErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_fetch_errors_total",
		Help: "The total number of errors during user fetch operations",
	})

	// UserUpdateCounter counts user update operations
	UserUpdateCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_update_total",
		Help: "The total number of user update operations",
	})

	// UserUpdateErrorCounter counts errors during user update operations
	UserUpdateErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_update_errors_total",
		Help: "The total number of errors during user update operations",
	})

	// UserDeleteCounter counts user delete operations
	UserDeleteCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_delete_total",
		Help: "The total number of user delete operations",
	})

	// UserDeleteErrorCounter counts errors during user delete operations
	UserDeleteErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "user_service_user_delete_errors_total",
		Help: "The total number of errors during user delete operations",
	})

	// RequestDurationHistogram measures request duration
	RequestDurationHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "user_service_request_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"handler", "status"})
)