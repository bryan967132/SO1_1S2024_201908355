package controller

import "github.com/gofiber/fiber/v2"

type Controller struct{}

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
