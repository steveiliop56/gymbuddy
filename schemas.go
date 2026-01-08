package main

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
	Name       string
	Excercises []ExerciseConfig
}

type ExerciseConfig struct {
	Name string
	Reps int
	Sets int
	Rest int // in seconds
}
