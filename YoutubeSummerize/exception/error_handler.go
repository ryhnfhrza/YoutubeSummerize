package exception

import (
	"net/http"

	"github.com/ryhnfhrza/YoutubeSummerize/helper"
	"github.com/ryhnfhrza/YoutubeSummerize/model/web"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				handleError(writer, request, err)
			}
		}()
		next.ServeHTTP(writer, request)
	})
}

func handleError(writer http.ResponseWriter , request *http.Request, err interface{}){
	if notFoundError(writer,request,err){
		return
	}

	if badRequestError(writer,request,err){
		return
	}


	internalServerError(writer,request,err)

}

func notFoundError(writer http.ResponseWriter , request *http.Request, err interface{}) bool {
	exception,ok := err.(NotFoundError)
	if ok{
		writer.Header().Set("Content-Type","application/json")
		writer.WriteHeader(http.StatusNotFound)

		webResponse := web.WebResponse{
			Code: http.StatusNotFound,
			Status: "NOT FOUND",
			Data: exception.Error,
		}
		helper.WriteToResponseBody(writer,webResponse)
		return true
	}else{
		return false
	}
}

func internalServerError(writer http.ResponseWriter , request *http.Request, err interface{}){
	writer.Header().Set("Content-Type","application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	webResponse := web.WebResponse{
		Code: http.StatusInternalServerError,
		Status: "INTERNAL SERVER ERROR",
		Data: err,
	}

	helper.WriteToResponseBody(writer,webResponse)
}

func badRequestError(writer http.ResponseWriter , request *http.Request, err interface{}) bool {
	exception,ok := err.(BadRequestError)
	if ok{
		writer.Header().Set("Content-Type","application/json")
		writer.WriteHeader(http.StatusBadRequest)

		webResponse := web.WebResponse{
			Code: http.StatusBadRequest,
			Status: "BAD REQUEST",
			Data: exception.Error,
		}
		helper.WriteToResponseBody(writer,webResponse)
		return true
	}else{
		return false
	}
}