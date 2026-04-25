package DBdrivers

import (
	"bytes"
	"context"
	"database/sql"
	"dbtool/DBinterface"
	"fmt"
	"os"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectToMySQL(conf DBinterface.Config) error {

	connString := DBinterface.CreateConnectionString(conf)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return fmt.Errorf("connection fail: %w", err)
	}
	defer db.Close()

	db.SetConnMaxLifetime(5 * time.Second)
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping fail: %w", err)
	}
	fmt.Println("Success!")
	return nil
}

func BackupMySQL(conf DBinterface.Config) error {
	if _, err := exec.LookPath("mysqldump"); err != nil {
		return fmt.Errorf("mysqldump not found: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []string{
		"-h", conf.Host,
		"-P", fmt.Sprintf("%d", conf.Port),
		"-u", conf.User,
		"--single-transaction",
		"--routines",
		"--triggers",
		conf.DBName,
	}

	cmd := exec.CommandContext(ctx, "mysqldump", args...)
	env := os.Environ()
	env = append(env, fmt.Sprintf("MYSQL_PWD=%s", conf.Password))
	cmd.Env = env

	var stderr bytes.Buffer
	outFile, err := os.Create(conf.FilePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	cmd.Stdout = outFile
	cmd.Stderr = &stderr

	fmt.Println("Starting backup")
	startTime := time.Now()
	err = cmd.Run()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("backup timed out")
		}
		errMsg := stderr.String()
		return fmt.Errorf("mysqldump failed: %w\nDetails: %s", err, errMsg)
	}

	duration := time.Since(startTime)
	fmt.Printf("Backup completed successfully in %v\nSaved to: %s\n", duration, conf.FilePath)
	return nil
}

func RestoreMySQL(conf DBinterface.Config) error {
	if _, err := exec.LookPath("mysql"); err != nil {
		return fmt.Errorf("mysql restore not found: %w", err)
	}
	if _, err := os.Stat(conf.FilePath); os.IsNotExist(err) {
		return fmt.Errorf("File not found: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	args := []string{
		"-h", conf.Host,
		"-P", fmt.Sprintf("%d", conf.Port),
		"-u", conf.User,
		conf.DBName,
		"--force",
	}
	cmd := exec.CommandContext(ctx, "mysql", args...)
	env := os.Environ()
	env = append(env, fmt.Sprintf("MYSQL_PWD=%s", conf.Password))
	cmd.Env = env

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	//fmt.Printf("WARNING!\n DB %s will be delete if it already exists!\n", conf.DBName)
	//time.Sleep(3 * time.Second)
	fmt.Println("Starting restore")
	startTime := time.Now()
	//if RecreateErr := Recreate(ctx, conf); RecreateErr != nil {
	//	return RecreateErr
	//}
	err := cmd.Run()

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("restore timed out")
		}
		errMsg := stderr.String()
		return fmt.Errorf("mysql restore failed: %w\nDetails: %s", err, errMsg)
	}

	duration := time.Since(startTime)
	fmt.Printf("Restore completed successfully in %v\n", duration)
	return nil
}
