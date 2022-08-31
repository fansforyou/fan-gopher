package rod

import (
	"errors"
	"fmt"
	"time"

	"github.com/fansforyou/fan-gopher/fans"
	"github.com/fansforyou/fan-gopher/model"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type RodGopher struct {
	browser *rod.Browser
}

func NewGopher(browser *rod.Browser) *RodGopher {
	return &RodGopher{
		browser: browser,
	}
}

func (rd *RodGopher) GetPostDetails(postID int64, creatorName string) (*model.Post, error) {
	page, err := rd.getPage(postID, creatorName)
	if err != nil {
		return nil, fmt.Errorf("unable to get page for post ID %d under creator '%s': %w", postID, creatorName, err)
	}

	actorNameAnchor, err := page.ElementX("//div[contains(@class, 'g-user-name')]")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve username anchor: %w", err)
	}
	actorName := actorNameAnchor.MustText()

	actorImage, err := xpath(page, "//a[contains(@class, 'g-avatar')]", "//img")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve actor image anchor: %w", err)
	}
	actorImageSrc := actorImage.MustAttribute("src")
	var actorImageURL string
	if actorImageSrc != nil {
		actorImageURL = *actorImageSrc
	}

	videoDescriptionDivs, err := page.ElementsX("//div[contains(@class, 'b-post__text')]")
	if err != nil {
		return nil, fmt.Errorf("failed to resolve the video description div: %w", err)
	}

	var videoDescription string
	if len(videoDescriptionDivs) > 0 {
		videoDescription = videoDescriptionDivs[0].MustText()
	}

	return &model.Post{
		Actors: []*model.ActorDetails{
			{
				ActorName:       actorName,
				ProfileImageURL: actorImageURL,
			},
		},
		VideoDetails: &model.VideoDetails{
			VideoDescription: videoDescription,
		},
	}, nil
}

func xpath(page *rod.Page, xpaths ...string) (*rod.Element, error) {
	var currentElement *rod.Element
	for xpathIndex, xpath := range xpaths {
		var xpathErr error
		if xpathIndex == 0 {
			currentElement, xpathErr = page.ElementX(xpath)
		} else {
			currentElement, xpathErr = currentElement.ElementX(xpath)
		}

		if xpathErr != nil {
			return nil, fmt.Errorf("unable to evaluate xpath '%s' (at index %d): %w", xpath, xpathIndex, xpathErr)
		}
	}
	return currentElement, nil
}

func (rd *RodGopher) VerifyExists(postID int64, creatorName string) (bool, error) {
	_, err := rd.getPage(postID, creatorName)
	if err != nil {
		var postNotFound *fans.PostNotFoundErr
		if errors.As(err, &postNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to get page for post %d under creator '%s': %w", postID, creatorName, err)
	}

	return true, nil
}

func (rd *RodGopher) getPage(postID int64, creatorName string) (*rod.Page, error) {
	requestPath := fmt.Sprintf("https://onlyfans.com/%d/%s", postID, creatorName)
	page, err := rd.browser.Page(proto.TargetCreateTarget{URL: requestPath})
	if err != nil {
		return nil, fmt.Errorf("failed to request page at '%s': %w", requestPath, err)
	}

	// Wait for the page to finish loading
	page.WaitRequestIdle(2*time.Second, nil, nil)()

	notFoundDivs, err := page.ElementsX("//div[contains(@class, 'b-404')]")
	if err != nil {
		return nil, fmt.Errorf("failed to search page at URL '%s' for 404-indicating divs: %w", requestPath, err)
	}

	if !notFoundDivs.Empty() {
		return nil, &fans.PostNotFoundErr{
			PostID:      postID,
			CreatorName: creatorName,
		}
	}

	return page, nil
}
