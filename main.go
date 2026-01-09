package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"strings"

	_ "embed"

	"github.com/goccy/go-yaml"
)

//go:embed base_page.html
var basePage string

//go:embed workout_page.html
var workoutPage string

//go:embed exercise_page.html
var exercisePage string

//go:embed styles.css
var stylesCSS string

func main() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	configEnv := os.Getenv("GYMBUDDY_CONFIG_PATH")
	configFlag := flag.String("config-path", "config.yml", "Path to the configuration file.")

	flag.Parse()

	configPath := configEnv

	if configPath == "" {
		configPath = *configFlag
	}

	content, err := os.ReadFile(configPath)

	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	var config AppConfig

	err = yaml.Unmarshal(content, &config)

	if err != nil {
		fmt.Println("Error parsing config file:", err)
		return
	}

	logger.Info("loaded config file")

	baseTmpl, err := template.New("base").Parse(basePage)

	if err != nil {
		logger.Error("error parsing template", "error", err)
		return
	}

	workoutTmpl, err := template.New("workout").Parse(workoutPage)

	if err != nil {
		logger.Error("error parsing template", "error", err)
		return
	}

	exerciseTmpl, err := template.New("exercise").Parse(exercisePage)

	if err != nil {
		logger.Error("error parsing template", "error", err)
		return
	}

	logger.Info("loaded templates")

	pageData := PageData{}

	for _, workout := range config.Workouts {
		id := strings.ReplaceAll(strings.ToLower(workout.Name), " ", "")
		workout.ID = id

		exercises := make([]ExerciseConfig, 0)

		for _, exercise := range workout.Exercises {
			exerciseID := strings.ReplaceAll(strings.ToLower(exercise.Name), " ", "")
			exercise.ID = exerciseID
			exercises = append(exercises, exercise)
		}

		workout.Exercises = exercises
		pageData.Workouts = append(pageData.Workouts, workout)
	}

	router := http.NewServeMux()

	router.HandleFunc("/", basePageHandler(logger, baseTmpl, pageData))
	router.HandleFunc("/workout/", workoutPageHandler(logger, workoutTmpl, pageData))
	router.HandleFunc("/exercise/", exercisePageHandler(logger, exerciseTmpl, pageData))
	router.HandleFunc("/styles.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write([]byte(stylesCSS))
	})
	router.HandleFunc("/healthz/", healthCheckHandler)

	logger.Info("starting server", "address", config.Address, "port", config.Port)

	err = http.ListenAndServe(fmt.Sprintf("%s:%d", config.Address, config.Port), router)

	if err != nil {
		logger.Error("server error", "error", err)
	}
}
