package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd representa el comando base cuando se llama sin ningún subcomando
var rootCmd = &cobra.Command{
	Use:   "prunus",
	Short: "Prunus - Business Management API",
	Long: `Prunus es una API REST diseñada para la gestión integral de negocios,
incluyendo empresas, sucursales, usuarios y productos.
Sigue el patrón de Arquitectura Limpia para asegurar escalabilidad y mantenimiento.`,
}

// Execute añade todos los comandos hijos al comando raíz y establece los flags apropiadamente.
// Esto es llamado por main.main(). Solo necesita ocurrir una vez para el rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Aquí se pueden definir flags globales.
	// El flag --config permite especificar un archivo de configuración diferente al .env predeterminado.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "archivo de configuración (por defecto es .env)")
}

// initConfig lee el archivo de configuración y las variables de entorno.
func initConfig() {
	if cfgFile != "" {
		// Usar el archivo de configuración del flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Buscar en el directorio actual el archivo ".env"
		viper.AddConfigPath(".")
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
	}

	// Leer variables de entorno que coincidan
	viper.AutomaticEnv()
	// Reemplazar puntos por guiones bajos para variables de entorno (ej. DB.HOST -> DB_HOST)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Si se encuentra un archivo de configuración, leerlo.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "✅ Usando archivo de configuración:", viper.ConfigFileUsed())
	} else {
		fmt.Fprintln(os.Stderr, "Aviso: No se pudo cargar el archivo .env, usando variables de entorno del sistema")
	}
}
