package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/viveknathani/numero/nlog"
	"github.com/viveknathani/numero/nparser"
)

// EvalRequest is the request body for the /api/v1/eval endpoint
type EvalRequest struct {
	Expression string            `json:"expression"`
	Variables  nparser.Variables `json:"variables,omitempty"`
}

// sendStandardResponse sends a standard response
func sendStandardResponse(
	c *fiber.Ctx,
	code int,
	data *map[string]interface{},
	message string,
) error {
	return c.Status(code).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}

// handle404 handles 404 errors
func handle404(c *fiber.Ctx) error {
	return sendStandardResponse(c, fiber.StatusNotFound, nil, "you seem lost!")
}

// Pool for reusing request objects
var evalPool = sync.Pool{
	New: func() interface{} {
		return new(EvalRequest)
	},
}

func main() {
	nlog.Info("hello from numero!")

	PORT := "8084"
	// Configure Fiber with optimized settings
	app := fiber.New(fiber.Config{
		Prefork:      true, // Uses multiple processes
		ServerHeader: "Numero",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		BodyLimit:    4 * 1024,   // 4KB - adjust based on your needs
		Concurrency:  256 * 1024, // Max concurrent connections
	})

	// Add middlewares
	app.Use(recover.New())
	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Post("/api/v1/eval", func(c *fiber.Ctx) error {
		// Get request object from pool
		req := evalPool.Get().(*EvalRequest)
		defer evalPool.Put(req)

		// Reset request fields
		req.Expression = ""
		req.Variables = nil

		if err := c.BodyParser(req); err != nil {
			return sendStandardResponse(c, fiber.StatusBadRequest, nil, err.Error())
		}

		np := nparser.New(req.Expression)
		for name, value := range req.Variables {
			np.SetVariable(name, value)
		}
		result, err := np.Run()
		if err != nil {
			return sendStandardResponse(c, fiber.StatusBadRequest, nil, err.Error())
		}
		return sendStandardResponse(c, fiber.StatusOK, &map[string]interface{}{
			"result": result,
		}, "success")
	})

	app.Use(handle404)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := app.Listen(":" + PORT)
		if err != nil {
			nlog.Error(err.Error())
			os.Exit(1)
		}
	}()

	nlog.Info("server is up 💯, url: http://localhost:" + PORT)

	<-done
	nlog.Info("goodbye 🙋")
}
