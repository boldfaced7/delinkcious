package link_manager

import (
	"errors"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
)

type LinkManager struct {
	linkStore          LinkStore
	socialGraphManager om.SocialGraphManager
	eventSink          om.LinkEvents
}

func (m *LinkManager) GetLinks(request om.GetLinksRequest) (result om.GetLinksResult, err error) {
	if request.Username == "" {
		err = errors.New("User name can't be empty")
		return
	}

	result, err = m.linkStore.GetLinks(request)
	return
}

func (m *LinkManager) AddLink(request om.AddLinkRequest) (err error) {
	if request.Url == "" {
		return errors.New("the URL can't be empty")
	}

	if request.Username == "" {
		return errors.New("the user name can't be empty")
	}

	link, err := m.linkStore.AddLink(request)
	if err != nil {
		return
	}

	followers, err := m.socialGraphManager.GetFollowers(request.Username)
	if err != nil {
		return
	}

	for follower, _ := range followers {
		m.eventSink.OnLinkAdded(follower, link)
	}

	return
}

func (m *LinkManager) UpdateLink(request om.UpdateLinkRequest) (err error) {
	if request.Url == "" {
		return errors.New("the URL can't be empty")
	}

	if request.Username == "" {
		return errors.New("the user name can't be empty")
	}

	link, err := m.linkStore.UpdateLink(request)
	if err != nil {
		return
	}

	followers, err := m.socialGraphManager.GetFollowers(request.Username)
	if err != nil {
		return
	}

	for follower, _ := range followers {
		m.eventSink.OnLinkUpdated(follower, link)
	}

	return
}

func (m *LinkManager) DeleteLink(username string, url string) (err error) {
	if url == "" {
		return errors.New("the URL can't be empty")
	}

	if username == "" {
		return errors.New("the user name can't be empty")
	}

	err = m.linkStore.DeleteLink(username, url)
	return
}

func NewLinkManager(linkStore LinkStore,
	socialGrpahManager om.SocialGraphManager,
	eventSink om.LinkEvents) (om.LinkManager, error) {
	if linkStore == nil {
		return nil, errors.New("link store")
	}

	if socialGrpahManager == nil {
		return nil, errors.New("social graph manager can't be nil")
	}

	if eventSink == nil {
		return nil, errors.New("event sink can't be nil")
	}

	return &LinkManager{
		linkStore:          linkStore,
		socialGraphManager: socialGrpahManager,
	}, nil
}