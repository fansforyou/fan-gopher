package rod_test

import (
	gopherrod "github.com/fansforyou/fan-gopher/fans/rod"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rod", Label("integration-test"), func() {
	var gopher *gopherrod.RodGopher

	BeforeEach(func() {
		browser := rod.New().Logger(utils.LoggerQuiet)
		connectErr := browser.Connect()
		Expect(connectErr).To(BeNil(), "opening the browser should not have failed")
		DeferCleanup(browser.MustClose)

		gopher = gopherrod.NewGopher(browser)
	})

	It("should verify that the post exists", func() {
		Expect(gopher.VerifyExists(387905677, "petittits")).To(BeTrue(), "the post should be reflected as existent")
	})

	It("should reflect that a post does not exist", func() {
		Expect(gopher.VerifyExists(387905678, "petittits")).To(BeFalse(), "the post should be reflected as non-existent")
	})

	It("should successfully retrieve a post's details", func() {
		postDetails, err := gopher.GetPostDetails(387905677, "petittits")
		Expect(err).To(BeNil(), "getting the post details should not fail")
		Expect(postDetails).ToNot(BeNil(), "post details should have been returned")

		Expect(postDetails.VideoDetails).ToNot(BeNil(), "video details should have been returned")
		Expect(postDetails.VideoDetails.VideoDescription).To(Equal("Slap it!🍑"), "the video description should be correct")

		Expect(postDetails.Actors).To(HaveLen(1), "a single actor should be listed")
		actor := postDetails.Actors[0]
		Expect(actor.ActorName).ToNot(BeEmpty(), "the actor name should be returned")
		Expect(actor.ProfileImageURL).ToNot(BeEmpty(), "the actor profile image should be returned")
	})
})
