package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gomarkdown/markdown"
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

	app := fiber.New(fiber.Config{
		Prefork:      true,
		ServerHeader: "numero",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
		BodyLimit:    4 * 1024,
		Concurrency:  256 * 1024,
	})

	app.Use(recover.New())

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		md, err := os.ReadFile("README.md")
		if err != nil {
			return sendStandardResponse(c, fiber.StatusInternalServerError, nil, "failed to read README.md")
		}

		html := markdown.ToHTML(md, nil, nil)

		template := `<!DOCTYPE html>
			<html>
			<head>
				<meta charset="utf-8">
				<meta name="viewport" content="width=device-width, initial-scale=1">
				<title>numero - Mathematical Expression Parser</title>
				<style>
					body { 
						max-width: 800px; 
						margin: 40px auto; 
						padding: 0 20px; 
						font-family: -apple-system, system-ui, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; 
						line-height: 1.5;
					}
					pre { 
						background: #f6f8fa; 
						padding: 16px; 
						border-radius: 6px; 
						overflow-x: auto;
					}
					code { font-family: monospace; }
					a { color: #0366d6; text-decoration: none; }
					a:hover { text-decoration: underline; }
				</style>
			</head>
			<body>%s</body></html>`
		styled := []byte(fmt.Sprintf(template, html))

		c.Set("Content-Type", "text/html")
		return c.Send(styled)
	})

	app.Post("/api/v1/eval", func(c *fiber.Ctx) error {
		req := evalPool.Get().(*EvalRequest)
		defer evalPool.Put(req)

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
