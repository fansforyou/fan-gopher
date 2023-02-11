package rod_test

import (
	"time"

	gopherrod "github.com/fansforyou/fan-gopher/fans/rod"
	"github.com/fansforyou/fan-gopher/model"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Rod", Label("integration-test"), func() {
	var gopher *gopherrod.RodGopher

	BeforeEach(func() {
		// Disable Leakless usage because Windows flags it as a virus
		rodLauncherURL := launcher.New().Leakless(false).MustLaunch()
		browser := rod.New().ControlURL(rodLauncherURL).Logger(utils.LoggerQuiet)
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
		Expect(actor.ActorName).To(Equal("Petittits"), "the actor name should be returned")
		Expect(actor.ProfileImageURL).ToNot(BeEmpty(), "the actor profile image should be returned")
	})

	Context("when anonymous access to posts is disallowed", func() {
		It("should successfully retrieve a post's details, but with a blank text description", func() {
			ticker := time.NewTimer(30 * time.Second)
			detailsChan := make(chan *model.Post)
			go func() {
				postDetails, err := gopher.GetPostDetails(380234283, "jocaramore")
				Expect(err).To(BeNil(), "getting the post details should not fail")
				detailsChan <- postDetails
			}()
			var postDetails *model.Post
			// Prepare to time out because the failure to read can, otherwise, cause an indefinite hang
			select {
			case <-ticker.C:
				panic("Test timed out waiting for post details")
			case p := <-detailsChan:
				postDetails = p
			}
			Expect(postDetails).ToNot(BeNil(), "post details should have been returned")

			Expect(postDetails.VideoDetails).ToNot(BeNil(), "video details should have been returned")
			Expect(postDetails.VideoDetails.VideoDescription).To(BeEmpty(), "because the description is unavailable, the description should be empty")

			Expect(postDetails.Actors).To(HaveLen(1), "a single actor should be listed")
			actor := postDetails.Actors[0]
			Expect(actor.ActorName).To(Equal("jocaramore"), "the actor name should be returned")
			Expect(actor.ProfileImageURL).ToNot(BeEmpty(), "the actor profile image should be returned")
		})
	})
})
