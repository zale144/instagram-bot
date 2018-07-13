package resource

import (
	"errors"
	"net/http"

	"github.com/zale144/instagram-bot/web/storage"

	"github.com/labstack/echo"
)

type JobResource struct{}

// get all jobs
func (jr JobResource) GetJobs(c echo.Context) error {
	jobs, err := storage.JobStorage{}.GetAll()
	if err != nil {
		err := errors.New("cannot get jobs")
		c.Error(echo.NewHTTPError(http.StatusBadRequest, err.Error()))
		return err
	}
	return c.JSON(http.StatusOK, jobs)
}
