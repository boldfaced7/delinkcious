package service

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/the-gigi/delinkcious/pb/news_service/pb"
	om "github.com/the-gigi/delinkcious/pkg/object_model"
	"net/http"
)

type getNewsRequest struct {
	Username   string `json:"username"`
	StartToken string `json:"start_token"`
}

type getNewsResult struct {
	Events    []*om.LinkManagerEvent `json:"events"`
	NextToken string                 `json:"next_token"`
	Err       string                 `json:"err"`
}

func decodeGetNewsRequest(_ context.Context, r *http.Request) (interface{}, error) { // http
	var request getNewsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

func encodeGetNewsResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func makeGetNewsEndpoint(svc om.NewsManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.GetNewsRequest)
		r, err := svc.GetNews(req)
		res := getNewsResult{}
		if err != nil {
			res.Err = err.Error()
		}
		return res, nil

		for _, e := range r.Events {
			res.Events = append(res.Events, e)
		}
		return res, nil
	}
}

/*

func newEvent(e *om.LinkManagerEvent) (event om.LinkManagerEvent) {
	event = om.LinkManagerEvent{
		EventType: (pb.EventType)(e.EventType),
		Username:  e.Username,
		Url:       e.Url,
	}

	seconds := e.Timestamp.Unix()
	nanos := (int32(e.Timestamp.UnixNano() - 1e9*seconds))

	return {}
	event.Timestamp = &timestamp.Timestamp{Seconds: seconds, Nanos: nanos}
	return
}

func makeGetNewsEndpoint(svc om.NewsManager) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(om.GetNewsRequest)
		r, err := svc.GetNews(req)
		res := &pb.GetNewsResponse{
			Events:    []*pb.Event{},
			NextToken: r.NextToken,
		}
		if err != nil {
			res.Err = err.Error()
		}
		for _, e := range r.Events {
			event := newEvent(e)
			res.Events = append(res.Events, event)
		}
		return res, nil
	}
}

func decodeGetNewsRequest(_ context.Context, r interface{}) (interface{}, error) {
	request := r.(*pb.GetNewsRequest)
	return om.GetNewsRequest{
		Username:   request.Username,
		StartToken: request.StartToken,
	}, nil
}

func encodeGetNewsResponse(_ context.Context, r interface{}) (interface{}, error) { // http
	return r, nil
}

*/

type handler struct {
	getNews grpctransport.Handler
}

func (s *handler) GetNews(ctx context.Context, r *pb.GetNewsRequest) (*pb.GetNewsResponse, error) {
	_, resp, err := s.getNews.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.GetNewsResponse), nil
}
