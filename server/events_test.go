package server_test

import (
	"fmt"
	"net/http"

	"github.com/nkcraddock/numzero/game"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("events integration tests", func() {
	var s *ServerHarness

	BeforeEach(func() {
		s = NewServerHarness()
		s.Authenticate("username", "password")
	})

	req_user := map[string]interface{}{"Name": "Chad"}
	req_coffee := map[string]interface{}{
		"code":   "coffee",
		"desc":   "Made a new pot of coffee",
		"points": 1,
	}

	Context("/events", func() {
		req_tooMuchCoffee := &game.Event{
			Player:      "shmurda",
			Description: "lunch break",
			Url:         "",
			Scores: []game.Score{
				game.Score{"coffee", 1000},
			},
		}

		It("stores an event for the player", func() {
			s.PUT("/players", &req_user)
			s.PUT("/rules", &req_coffee)

			res := s.POST("/events", req_tooMuchCoffee)
			Ω(res.Code).Should(Equal(http.StatusCreated))
		})
	})

	Context("POST /events", func() {
		It("adds a new event", func() {
			res := s.POST("/events", &req_coffee)
			Ω(res.Code).Should(Equal(http.StatusCreated))
		})
	})

	Context("GET /events", func() {
		It("retrieves a event", func() {
			res := s.POST("/events", &req_coffee)
			created := new(game.Event)
			err := s.Parse(res, created)
			Ω(err).ShouldNot(HaveOccurred())

			event := new(game.Event)
			res = s.GET(fmt.Sprintf("/events/%s", created.Id), event)
			Ω(res.Code).Should(Equal(http.StatusOK))
			Ω(event.Description).Should(Equal(req_coffee["desc"]))
		})
	})
})
