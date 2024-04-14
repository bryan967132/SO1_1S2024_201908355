package controller

import (
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

type Graph struct {
	Transitions       fiber.Map
	ThereIsAnyProcess bool
	CurrentPID        int
}

type Controller struct {
	DB    *sql.DB
	err   error
	Graph *Graph
}

func NewController() *Controller {
	ctrl := &Controller{}
	ctrl.connect()
	ctrl.Graph = &Graph{Transitions: fiber.Map{}, ThereIsAnyProcess: false}
	return ctrl
}

func (c *Controller) connect() bool {
	c.DB, c.err = sql.Open("mysql", "root:mysqlpass@tcp(mysql:3306)/P1SO1")
	if c.err != nil {
		return false
	}

	c.err = c.DB.Ping()
	if c.err != nil {
		return false
	}

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
		fmt.Println("ERROR EN cat /proc/ram_cpu", err.Error())
		return ctx.JSON(fiber.Map{
			"status":      "Error 1",
			"descripcion": err.Error(),
		})
	}

	var datos fiber.Map

	err = json.Unmarshal(stdout.Bytes(), &datos)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":      "Error 2",
			"descripcion": err.Error(),
		})
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

	cmd = exec.Command("ps", "-eo", "%cpu")
	output, err := cmd.Output()
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":      "Error 3",
			"descripcion": err.Error(),
		})
	}

	lines := strings.Split(string(output), "\n")
	totalCPU := 0.0
	for _, line := range lines[1:] {
		cpuUsageStr := strings.TrimSpace(line)
		if cpuUsageStr == "" {
			continue
		}
		cpuUsage, err := strconv.ParseFloat(cpuUsageStr, 64)
		if err != nil {
			continue
		}
		totalCPU += cpuUsage
	}

	if totalCPU > 100 {
		totalCPU = 98.89
	}

	return ctx.JSON(fiber.Map{
		"status":         "Success",
		"percentusedRAM": (totalram - freeram) * 100 / totalram,
		"percentusedCPU": totalCPU,
	})
}

func (c *Controller) InsRAMCPU(ctx *fiber.Ctx) error {
	var reqBody Reg
	if err := ctx.BodyParser(&reqBody); err != nil {
		return ctx.JSON(fiber.Map{
			"status":      "Query RAM CPU error 1",
			"descripcion": err.Error(),
		})
	}

	_, err := c.DB.Exec("INSERT INTO RAM (usado, disponible, tiempo) VALUE (?, ?, ?)", reqBody.UsadoRAM, reqBody.DisponibleRAM, time.Now())
	if err != nil {
		fmt.Println("ERROR EN INSERT RAM", err.Error())
		return ctx.JSON(fiber.Map{
			"status":      "Query INSERT RAM error 2",
			"descripcion": err.Error(),
		})
	}

	_, err = c.DB.Exec("INSERT INTO CPU (usado, disponible, tiempo) VALUE (?, ?, ?)", reqBody.UsadoCPU, reqBody.DisponibleCPU, time.Now())
	if err != nil {
		fmt.Println("ERROR EN INSERT CPU", err.Error())
		return ctx.JSON(fiber.Map{
			"status":      "Query INSERT CPU error 3",
			"descripcion": err.Error(),
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
		fmt.Println("ERROR EN SELECT RAM", err.Error())
		return ctx.JSON(fiber.Map{
			"status":      "Query RAM error 1",
			"descripcion": err.Error(),
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
				"status":      "Query RAM error 2",
				"descripcion": err.Error(),
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
		fmt.Println("ERROR EN SELECT CPU", err.Error())
		return ctx.JSON(fiber.Map{
			"status":      "Query CPU error 1",
			"descripcion": err.Error(),
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
				"status":      "Query CPU error 2",
				"descripcion": err.Error(),
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
			"status":      "Error Pids 1",
			"descripcion": err.Error(),
		})
	}

	var datos fiber.Map

	err = json.Unmarshal(stdout.Bytes(), &datos)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":      "Error Pids 2",
			"descripcion": err.Error(),
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
			"status":      "Error Proc 1",
			"descripcion": err.Error(),
		})
	}

	var datos fiber.Map

	err = json.Unmarshal(stdout.Bytes(), &datos)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":      "Error Proc 2",
			"descripcion": err.Error(),
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

func (c *Controller) ThereIsProc(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": c.Graph.ThereIsAnyProcess,
		"PID":    c.Graph.CurrentPID,
		"graph":  c.Graph.Transitions,
	})
}

func (c *Controller) Start(ctx *fiber.Ctx) error {
	if !c.Graph.ThereIsAnyProcess {
		cmd := exec.Command("sleep", "infinity")
		err := cmd.Start()
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":      "Fail Create",
				"descripcion": err.Error(),
			})
		}

		c.Graph.ThereIsAnyProcess = true
		c.Graph.CurrentPID = cmd.Process.Pid

		c.Graph.Transitions = fiber.Map{
			"new": fiber.Map{
				"status": "Nothing",
				"to":     []string{"ready"},
			},
			"ready": fiber.Map{
				"status": "Nothing",
				"to":     []string{"running"},
			},
			"running": fiber.Map{
				"status": "Current",
				"to":     []string{},
			},
			"terminated": fiber.Map{
				"status": "Nothing",
				"to":     []string{},
			},
		}

		return ctx.JSON(fiber.Map{
			"status": "Success",
			"PID":    cmd.Process.Pid,
			"graph":  c.Graph.Transitions,
		})
	}
	return ctx.JSON(fiber.Map{
		"status": "Current Process",
		"PID":    c.Graph.CurrentPID,
		"graph":  c.Graph.Transitions,
	})
}

func (c *Controller) contain(arr []string, value string) bool {
	for _, el := range arr {
		if el == value {
			return true
		}
	}
	return false
}

func (c *Controller) Stop(ctx *fiber.Ctx) error {
	if c.Graph.ThereIsAnyProcess {
		pidStr := ctx.Params("pid")
		if pidStr == "" {
			return ctx.JSON(fiber.Map{
				"status": "PID Invalid",
			})
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":      "PID Invalid",
				"descripcion": err.Error(),
			})
		}

		cmd := exec.Command("kill", "-SIGSTOP", strconv.Itoa(pid))
		err = cmd.Run()
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":      "Non-exist Process",
				"PID":         pid,
				"descripcion": err.Error(),
			})
		}

		c.Graph.Transitions["running"].(fiber.Map)["status"] = "Nothing"
		if !c.contain(c.Graph.Transitions["running"].(fiber.Map)["to"].([]string), "ready") {
			slice := c.Graph.Transitions["running"].(fiber.Map)["to"].([]string)
			slice = append(slice, "ready")
			c.Graph.Transitions["running"].(fiber.Map)["to"] = slice
		}
		c.Graph.Transitions["ready"].(fiber.Map)["status"] = "Current"

		return ctx.JSON(fiber.Map{
			"status": "Success",
			"graph":  c.Graph.Transitions,
		})
	}
	return ctx.JSON(fiber.Map{
		"status": "Non-exist Process Current",
	})
}

func (c *Controller) Resume(ctx *fiber.Ctx) error {
	if c.Graph.ThereIsAnyProcess {
		pidStr := ctx.Params("pid")
		if pidStr == "" {
			return ctx.JSON(fiber.Map{
				"status": "PID Invalid",
			})
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":      "PID Invalid",
				"descripcion": err.Error(),
			})
		}

		cmd := exec.Command("kill", "-SIGCONT", strconv.Itoa(pid))
		err = cmd.Run()
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":      "Non-exist Process",
				"PID":         pid,
				"descripcion": err.Error(),
			})
		}

		c.Graph.Transitions["ready"].(fiber.Map)["status"] = "Nothing"
		c.Graph.Transitions["running"].(fiber.Map)["status"] = "Current"

		return ctx.JSON(fiber.Map{
			"status": "Success",
			"graph":  c.Graph.Transitions,
		})
	}
	return ctx.JSON(fiber.Map{
		"status": "Non-exist Process Current",
	})
}

func (c *Controller) Kill(ctx *fiber.Ctx) error {
	if c.Graph.ThereIsAnyProcess {
		pidStr := ctx.Params("pid")
		if pidStr == "" {
			return ctx.JSON(fiber.Map{
				"status": "PID Invalid",
			})
		}

		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":      "PID Invalid",
				"descripcion": err.Error(),
			})
		}

		cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
		err = cmd.Run()
		if err != nil {
			return ctx.JSON(fiber.Map{
				"status":      "Non-exist Process",
				"PID":         pid,
				"descripcion": err.Error(),
			})
		}

		c.Graph.ThereIsAnyProcess = false

		if c.Graph.Transitions["ready"].(fiber.Map)["status"] == "Current" {
			c.Graph.Transitions["ready"].(fiber.Map)["status"] = "Nothing"
			if !c.contain(c.Graph.Transitions["ready"].(fiber.Map)["to"].([]string), "terminated") {
				slice := c.Graph.Transitions["ready"].(fiber.Map)["to"].([]string)
				slice = append(slice, "terminated")
				c.Graph.Transitions["ready"].(fiber.Map)["to"] = slice
			}
		}
		if c.Graph.Transitions["running"].(fiber.Map)["status"] == "Current" {
			c.Graph.Transitions["running"].(fiber.Map)["status"] = "Nothing"
			if !c.contain(c.Graph.Transitions["running"].(fiber.Map)["to"].([]string), "terminated") {
				slice := c.Graph.Transitions["running"].(fiber.Map)["to"].([]string)
				slice = append(slice, "terminated")
				c.Graph.Transitions["running"].(fiber.Map)["to"] = slice
			}
		}
		c.Graph.Transitions["terminated"].(fiber.Map)["status"] = "Current"

		return ctx.JSON(fiber.Map{
			"status": "Success",
			"graph":  c.Graph.Transitions,
		})
	}
	return ctx.JSON(fiber.Map{
		"status": "Non-exist Process Current",
	})
}
