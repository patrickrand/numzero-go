package server_test

import (
	"net/http"

	"github.com/nkcraddock/numzero"
	"github.com/nkcraddock/numzero/game"
	"github.com/nkcraddock/numzero/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("players integration tests", func() {
	var s *ServerHarness

	BeforeEach(func() {
		authStore := numzero.NewMemoryStore()
		store := game.NewMemoryStore()
		store.SaveRule(game.Rule{"coffee", "made coffee", 1})
		store.SaveRule(game.Rule{"highfive", "high-fived someone", -10})

		s = NewServerHarness(authStore, store)
		s.Authenticate("username", "password")
	})

	req_chad := map[string]interface{}{"Name": "Chad"}
	req_roger := map[string]interface{}{"Name": "Roger"}

	Context("PUT /players", func() {
		It("adds a new player", func() {
			res := s.PUT("/players", &req_chad)
			Ω(res.Code).Should(Equal(http.StatusCreated))
		})

		It("updates an existing payer", func() {
			s.PUT("/players", &req_chad)

			req_update := map[string]interface{}{
				"Name":  "Chad",
				"Score": 1000,
			}

			s.PUT("/players", &req_update)

			p := &game.Player{}
			s.GET("/players/Chad", p)
			Ω(p.Score).Should(Equal(1000))
		})
	})

	Context("GET /players", func() {
		It("gets a player", func() {
			s.PUT("/players", &req_chad)

			p := &game.Player{}
			res := s.GET("/players/Chad", p)

			Ω(res.Code).Should(Equal(http.StatusOK))
			Ω(p.Name).Should(Equal(req_chad["Name"]))
		})

		It("gets a list of players", func() {
			s.PUT("/players", &req_chad)
			s.PUT("/players", &req_roger)

			var players []game.Player

			res := s.GET("/players", &players)
			Ω(res.Code).Should(Equal(http.StatusOK))
			Ω(players).Should(HaveLen(2))
		})
	})

	Context("Activities", func() {
		req_activity := map[string]interface{}{
			"desc": "Breakroom Visit - 5/1/2015 15:15 EST",
			"scores": map[string]int{
				"coffee":   2,
				"highfive": 1,
			},
		}

		Context("POST /players/{name}/activities", func() {

			It("adds an activity for a player", func() {
				s.PUT("/players", &req_chad)
				res := s.POST("/players/Chad/activities", req_activity)
				Ω(res.Code).Should(Equal(http.StatusOK))
			})
		})

		Context("GET /players/{name}/activities", func() {
			It("Gets a list of activities for a player", func() {
				s.PUT("/players", &req_chad)
				s.POST("/players/Chad/activities", req_activity)

				var act []server.Activity
				res := s.GET("/players/Chad/activities", &act)
				Ω(res.Code).Should(Equal(http.StatusOK))
				Ω(act).Should(HaveLen(1))
				Ω(act[0].Scores).Should(HaveLen(2))
			})
		})
	})
})