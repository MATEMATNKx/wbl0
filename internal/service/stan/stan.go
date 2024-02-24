package stan

import (
	"l0/internal/config"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/nats-io/stan.go/pb"
)

type OrderService interface {
	Create(orderUID, data string)
	Get(orderUID string) (string, error)
}
type Service struct {
	orderSvc OrderService

	NatsConnect *nats.Conn
	StanConnect stan.Conn
	Sub         stan.Subscription

	url       string
	clusterID string
	clientID  string
	subject   string
}

func New(cfg *config.Stan, orderSvc OrderService) *Service {
	return &Service{
		orderSvc:  orderSvc,
		url:       cfg.URL,
		clusterID: cfg.ClusterID,
		clientID:  cfg.ClientID,
		subject:   cfg.Subject,
	}
}
func (s *Service) Start() {
	nc, err := nats.Connect(s.url, nats.Name("Orders reader"))
	if err != nil {
		log.Fatal(err)
	}
	sc, err := stan.Connect(
		s.clusterID, s.clientID, stan.NatsConn(nc),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)

		}),
	)
	if err != nil {
		log.Fatalf("[stan.go] Can't connect: %v. \nMake sure a Nats Streaming Server is running at: %s", err, s.url)

	}
	log.Printf("[stan.go] Connected to %s clusterID: [%s] clientID: [%s]\n", s.url, s.clusterID, s.clientID)

	sub, err := sc.QueueSubscribe(
		s.subject,
		"",
		s.handleMessage,
		stan.StartAt(pb.StartPosition_NewOnly),
		stan.DurableName(""),
	)
	if err != nil {
		sc.Close()
		log.Fatal(err)
	}
	log.Printf("[stan.go] Listening on [%s], cliendID=[%s]\n", s.subject, s.clientID)
	s.Sub = sub
	s.NatsConnect = nc
	s.StanConnect = sc
}
func (s *Service) Stop() {
	s.StanConnect.Close()
	s.NatsConnect.Close()
}
