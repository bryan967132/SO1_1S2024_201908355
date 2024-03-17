package controller

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
)

type Reg struct {
	UsadoRAM      float64 `json:"usadoram"`
	DisponibleRAM float64 `json:"disponibleram"`
	UsadoCPU      float64 `json:"usadocpu"`
	DisponibleCPU float64 `json:"disponiblecpu"`
}

type Controller struct {
	DB  *sql.DB
	err error
}

func NewController() *Controller {
	ctrl := &Controller{}
	ctrl.connect()
	return ctrl
}

func (c *Controller) connect() bool {
	c.DB, c.err = sql.Open("mysql", "root:mysqlpass@tcp(127.0.0.1:3306)/P1SO1")
	if c.err != nil {
		return false
	}

	c.err = c.DB.Ping()
	if c.err != nil {
		return false
	}

	c.DB.SetMaxOpenConns(100 * 1024)
	return true
}

func (c *Controller) Running(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "Server is running!!!",
	})
}

func (c *Controller) Cpuram(ctx *fiber.Ctx) error {
	cmd := exec.Command("cat", "/proc/ram_cpu")
	// Crear un buffer para capturar la salida estándar y los errores
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Ejecutar el comando
	err := cmd.Run()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Error 1",
		})
	}

	cmd_ := exec.Command("mpstat")
	var stdout_, stderr_ bytes.Buffer
	cmd_.Stdout = &stdout_
	cmd_.Stderr = &stderr_

	err = cmd_.Run()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Error 2",
		})
	}

	var datos fiber.Map

	err = json.Unmarshal(stdout.Bytes(), &datos)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Error 3",
		})
	}

	dataCPU := []float64{}

	scanner := bufio.NewScanner(strings.NewReader(stdout_.String()))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "all") {
			for _, d := range strings.Fields(line) {
				n, _ := strconv.ParseFloat(d, 64)
				dataCPU = append(dataCPU, n)
			}
		}
	}

	var totalram int64
	switch v := datos["totalram"].(type) {
	case int64:
		totalram = v
	case float64:
		totalram = int64(v)
	}

	var freeram int64
	switch v := datos["freeram"].(type) {
	case int64:
		freeram = v
	case float64:
		freeram = int64(v)
	}

	return ctx.JSON(fiber.Map{
		"status":         "Success",
		"percentusedRAM": (totalram - freeram) * 100 / totalram,
		"percentusedCPU": dataCPU[11] * 100 / (dataCPU[2] + dataCPU[3] + dataCPU[4] + dataCPU[5] + dataCPU[6] + dataCPU[7] + dataCPU[8] + dataCPU[9] + dataCPU[10] + dataCPU[11]),
	})
}

func (c *Controller) InsRAMCPU(ctx *fiber.Ctx) error {
	var reqBody Reg
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Query RAM CPU error 1",
		})
	}

	_, err := c.DB.Exec("INSERT INTO RAM (usado, disponible, tiempo) VALUE (?, ?, ?)", reqBody.UsadoRAM, reqBody.DisponibleRAM, time.Now())
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Query INSERT RAM error 2",
		})
	}

	_, err = c.DB.Exec("INSERT INTO CPU (usado, disponible, tiempo) VALUE (?, ?, ?)", reqBody.UsadoCPU, reqBody.DisponibleCPU, time.Now())
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Query INSERT CPU error 3",
		})
	}

	return ctx.JSON(fiber.Map{
		"status": "Registros RAM CPU insertados.",
	})
}

func (c *Controller) History(ctx *fiber.Ctx) error {
	query := `SELECT * FROM RAM;`
	rows, err := c.DB.Query(query)
	if err != nil {
		defer rows.Close()
		return ctx.JSON(fiber.Map{
			"status": "Query RAM error 1",
		})
	}
	defer rows.Close()

	responseRAM := []fiber.Map{}

	for rows.Next() {
		var id int
		var usado float64
		var disponible float64
		var tiempo string
		if err := rows.Scan(&id, &usado, &disponible, &tiempo); err != nil {
			return ctx.JSON(fiber.Map{
				"status": "Query RAM error 2",
			})
		}
		responseRAM = append(responseRAM, fiber.Map{
			"id":         id,
			"usado":      usado,
			"disponible": disponible,
			"tiempo":     tiempo,
		})
	}

	query = `SELECT * FROM CPU;`
	rows, err = c.DB.Query(query)
	if err != nil {
		defer rows.Close()
		return ctx.JSON(fiber.Map{
			"status": "Query CPU error 1",
		})
	}
	defer rows.Close()

	responseCPU := []fiber.Map{}

	for rows.Next() {
		var id int
		var usado float64
		var disponible float64
		var tiempo string
		if err := rows.Scan(&id, &usado, &disponible, &tiempo); err != nil {
			return ctx.JSON(fiber.Map{
				"status": "Query CPU error 2",
			})
		}
		responseCPU = append(responseCPU, fiber.Map{
			"id":         id,
			"usado":      usado,
			"disponible": disponible,
			"tiempo":     tiempo,
		})
	}

	return ctx.JSON(fiber.Map{
		"status": "Success",
		"ram":    responseRAM,
		"cpu":    responseCPU,
	})
}

func (c *Controller) Pids(ctx *fiber.Ctx) error {
	cmd := exec.Command("cat", "/proc/ram_cpu")
	// Crear un buffer para capturar la salida estándar y los errores
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Ejecutar el comando
	err := cmd.Run()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Error Pids 1",
		})
	}

	var datos fiber.Map

	err = json.Unmarshal(stdout.Bytes(), &datos)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Error 3",
		})
	}

	pids := []float64{}

	for _, p := range datos["processes"].([]interface{}) {
		if convertedMap, ok := p.(map[string]interface{}); ok {
			if v, o := fiber.Map(convertedMap)["pid"].(float64); o {
				pids = append(pids, v)
			}
		}
	}

	return ctx.JSON(fiber.Map{
		"status": "Success",
		"pids":   pids,
	})
}

func (c *Controller) Proc(ctx *fiber.Ctx) error {
	pid := ctx.Params("pid")

	cmd := exec.Command("cat", "/proc/ram_cpu")
	// Crear un buffer para capturar la salida estándar y los errores
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Ejecutar el comando
	err := cmd.Run()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Error Pids 1",
		})
	}

	var datos fiber.Map

	err = json.Unmarshal(stdout.Bytes(), &datos)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status": "Error 3",
		})
	}

	for _, p := range datos["processes"].([]interface{}) {
		if convertedMap, ok := p.(map[string]interface{}); ok {
			if fmt.Sprint(pid) == fmt.Sprint(fiber.Map(convertedMap)["pid"]) {
				return ctx.JSON(fiber.Map{
					"status": "Success",
					"proc":   fiber.Map(convertedMap),
				})
			}
		}
	}

	return ctx.JSON(fiber.Map{
		"status": "Success",
		"proc":   "proc",
	})
}
