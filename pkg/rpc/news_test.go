package rpc_test

import (
	"testing"

	"apisrv/pkg/db/test"
	"apisrv/pkg/newsportal"
	"apisrv/pkg/rpc"
	. "github.com/smartystreets/goconvey/convey"
)

func initRPC(t *testing.T) *rpc.NewsService {
	db, _ := test.Setup(t)
	service := newsportal.NewNewsService(db)
	srv := rpc.NewNewsService(service)

	So(srv, ShouldNotBeNil)

	return srv
}

func TestDB_NewsService_Get(t *testing.T) {
	Convey("Test NewsService Get", t, func() {
		ctx := t.Context()
		srv := initRPC(t)

		Convey("List and count", func() {
			positiveCases := []struct {
				Name        string
				Req         rpc.NewsListReq
				LenExpected int
			}{
				{
					Name:        "Without filters",
					LenExpected: 1,
				},
				{
					Name:        "With category filter",
					Req:         rpc.NewsListReq{CategoryID: 1},
					LenExpected: 1,
				},
				{
					Name:        "With tag filter",
					Req:         rpc.NewsListReq{TagID: 1},
					LenExpected: 1,
				},
				{
					Name:        "With category & tag filter",
					Req:         rpc.NewsListReq{CategoryID: 1, TagID: 1},
					LenExpected: 1,
				},
				{
					Name:        "With unknown category",
					Req:         rpc.NewsListReq{CategoryID: 100500},
					LenExpected: 0,
				},
				{
					Name:        "With unknown tag",
					Req:         rpc.NewsListReq{TagID: 100500},
					LenExpected: 0,
				},
			}

			for _, testCase := range positiveCases {
				Convey(testCase.Name, func() {
					list, err := srv.Get(ctx, testCase.Req)

					So(err, ShouldBeNil)
					So(list, ShouldNotBeNil)
					So(list, ShouldHaveLength, testCase.LenExpected)

					count, err := srv.Count(ctx, testCase.Req)

					So(err, ShouldBeNil)
					So(count, ShouldEqual, testCase.LenExpected)
				})
			}
		})
	})
}

func TestDB_NewsService_GetByID(t *testing.T) {
	Convey("Test NewsService GetByID", t, func() {
		ctx := t.Context()
		srv := initRPC(t)

		Convey("Valid ID", func() {
			news, err := srv.GetByID(ctx, 5)

			So(err, ShouldBeNil)
			So(news, ShouldNotBeNil)
			So(news.ID, ShouldEqual, 5)
		})

		Convey("Unknown ID", func() {
			news, err := srv.GetByID(ctx, 100500)

			So(err, ShouldBeError)
			So(news, ShouldBeNil)
		})
	})
}
