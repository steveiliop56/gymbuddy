package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func basePageHandler(logger *slog.Logger, tmpl *template.Template, pageData PageData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, pageData)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			logger.Error("error rendering template", "error", err)
			return
		}
	}
}

func workoutPageHandler(logger *slog.Logger, tmpl *template.Template, pageData PageData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		param := strings.TrimPrefix(r.URL.Path, "/workout/")

		if param == "" {
			http.Error(w, "Workout not specified", http.StatusBadRequest)
			return
		}

		// yes it is inefficient but workouts are few and this is simple
		var workout *WorkoutConfig

		for _, wkt := range pageData.Workouts {
			if wkt.ID == param {
				workout = &wkt
				break
			}
		}

		if workout == nil {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}

		err := tmpl.Execute(w, *workout)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			logger.Error("error rendering template", "error", err)
			return
		}
	}
}

func exercisePageHandler(logger *slog.Logger, tmpl *template.Template, pageData PageData) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := strings.TrimPrefix(r.URL.Path, "/exercise/")

		workoutParam, exerciseParam, ok := strings.Cut(params, "/")

		if !ok {
			http.Error(w, "Workout or exercise not specified", http.StatusBadRequest)
			return
		}

		// yes it is inefficient but workouts are few and this is simple
		var workout *WorkoutConfig

		for _, wkt := range pageData.Workouts {
			if wkt.ID == workoutParam {
				workout = &wkt
				break
			}
		}

		if workout == nil {
			http.Error(w, "Workout not found", http.StatusNotFound)
			return
		}

		// yeah we get it, inefficient, simple
		var exercise *ExerciseConfig

		for _, ex := range workout.Exercises {
			if ex.ID == exerciseParam {
				exercise = &ex
				break
			}
		}

		if exercise == nil {
			http.Error(w, "Exercise not found", http.StatusNotFound)
			return
		}

		err := tmpl.Execute(w, *exercise)

		if err != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
			logger.Error("error rendering template", "error", err)
			return
		}
	}
}
