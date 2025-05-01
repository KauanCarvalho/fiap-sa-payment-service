package middleware_test

import (
	"context"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var ctx context.Context

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	ctx = context.Background()

	os.Exit(m.Run())
}
