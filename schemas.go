package main

// App config

type AppConfig struct {
	Port     int
	Address  string
	Auth     AuthConfig
	Workouts []WorkoutConfig
}

type AuthConfig struct {
	Username string
	Password string
}

type WorkoutConfig struct {
	ID        string `yaml:"-"`
	Name      string
	Exercises []ExerciseConfig
}

type ExerciseConfig struct {
	ID   string `yaml:"-"`
	Name string
	Reps int
	Sets int
	Rest int // in seconds
}

// Page data

type PageData struct {
	Workouts []WorkoutConfig
}
