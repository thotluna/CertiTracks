package testcontainer

type SilentLogger struct{}

func (l *SilentLogger) Printf(format string, args ...interface{}) {
	// No-op
}
func (l *SilentLogger) Fatalf(format string, args ...interface{}) {
	// testcontainers-go espera este método. Si ocurre un fatal, por lo general
	// indica un problema grave en el entorno de Docker que debería detener los tests.
	// Aquí puedes decidir si quieres que realmente detenga la ejecución o simplemente loguee.
	// Para el silencio total, también podría ser No-op, pero con cautela.
	// log.Fatalf(format, args...) // Esto detendría el test
}
