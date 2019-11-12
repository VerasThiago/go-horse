package middleware

import (
	"github.com/labbsr0x/go-horse/filters"
	"github.com/labbsr0x/go-horse/util"
	"github.com/kataras/iris/context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)


const (
	RequestBodyKey  = "requestBody"
	ResponseBodyKey = "responseBody"
)



func ResquestFilter(filter * filters.FilterManager) context.Handler {
	return func(ctx context.Context) {
		util.SetFilterContextValues(ctx)

		if ctx.Request().Body != nil {
			requestBody, err := ioutil.ReadAll(ctx.Request().Body)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"request": ctx.String(),
					"error": err.Error(),
				}).Errorf("Error parsing request body in the middleware")
			}
			ctx.Values().Set(RequestBodyKey, string(requestBody))
		}

		ctx.Values().Set("path", ctx.Request().URL.Path)

		result , err := filter.RunRequestFilters(ctx, RequestBodyKey)

		writer := ctx.ResponseWriter()
		ctx.ResetResponseWriter(writer)

		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Errorf("Error during the execution of REQUEST filters")
			writer.WriteString(err.Error())
			ctx.StopExecution()
			return
		}

		if result.Next{
			ctx.Next()
		}
	}
}