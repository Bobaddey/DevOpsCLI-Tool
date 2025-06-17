package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	tmpl "text/template"

	"github.com/spf13/cobra"
)

// Config holds the CLI configuration
type Config struct {
	TerraformPath string
	WorkspaceDir  string
	GitRemote     string
	GitBranch     string
}

// PipelineTemplate represents a CI/CD pipeline template
type PipelineTemplate struct {
	Name        string
	Language    string
	BuildSteps  []string
	TestSteps   []string
	DeploySteps []string
}

var config Config

func main() {
	// Initialize default config with hardcoded Terraform path
	config = Config{
		TerraformPath: "terraform",
	}

	// Prompt for Workspace directory
	fmt.Print("Enter workspace directory (press Enter for default './workspace'): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	workspaceDir := strings.TrimSpace(scanner.Text())
	if workspaceDir == "" {
		workspaceDir = "./workspace"
	}
	config.WorkspaceDir = workspaceDir

	// Prompt for Git branch
	fmt.Print("Enter Git branch (press Enter for default 'main'): ")
	scanner.Scan()
	gitBranch := strings.TrimSpace(scanner.Text())
	if gitBranch == "" {
		gitBranch = "main"
	}
	config.GitBranch = gitBranch

	// Create workspace directory if it doesn't exist
	if err := os.MkdirAll(config.WorkspaceDir, 0755); err != nil {
		log.Fatalf("Failed to create workspace directory: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "devopsctl",
		Short: "A DevOps automation CLI tool",
		Long:  "A comprehensive CLI tool for running Terraform, bash scripts, and managing CI/CD pipelines",
	}

	// Add subcommands
	rootCmd.AddCommand(terraformCmd())
	rootCmd.AddCommand(scriptCmd())
	rootCmd.AddCommand(pipelineCmd())
	rootCmd.AddCommand(configCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Terraform commands
func terraformCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "terraform",
		Short:   "Run Terraform operations",
		Aliases: []string{"tf"},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "init [directory]",
		Short: "Initialize Terraform in specified directory",
		Args:  cobra.MaximumNArgs(1),
		Run:   runTerraformInit,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "plan [directory]",
		Short: "Run Terraform plan",
		Args:  cobra.MaximumNArgs(1),
		Run:   runTerraformPlan,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "apply [directory]",
		Short: "Run Terraform apply",
		Args:  cobra.MaximumNArgs(1),
		Run:   runTerraformApply,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "destroy [directory]",
		Short: "Run Terraform destroy",
		Args:  cobra.MaximumNArgs(1),
		Run:   runTerraformDestroy,
	})

	return cmd
}

func runTerraformInit(cmd *cobra.Command, args []string) {
	dir := getCurrentDir(args)
	runTerraformCommand(dir, "init")
}

func runTerraformPlan(cmd *cobra.Command, args []string) {
	dir := getCurrentDir(args)
	runTerraformCommand(dir, "plan")
}

func runTerraformApply(cmd *cobra.Command, args []string) {
	dir := getCurrentDir(args)
	fmt.Print("Are you sure you want to apply these changes? (yes/no): ")
	if !confirmAction() {
		fmt.Println("Apply cancelled.")
		return
	}
	runTerraformCommand(dir, "apply", "-auto-approve")
}

func runTerraformDestroy(cmd *cobra.Command, args []string) {
	dir := getCurrentDir(args)
	fmt.Print("Are you sure you want to destroy these resources? (yes/no): ")
	if !confirmAction() {
		fmt.Println("Destroy cancelled.")
		return
	}
	runTerraformCommand(dir, "destroy", "-auto-approve")
}

func runTerraformCommand(dir string, args ...string) {
	cmd := exec.Command(config.TerraformPath, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running: terraform %s in %s\n", strings.Join(args, " "), dir)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Terraform command failed: %v", err)
	}
}

// Script commands
func scriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "script",
		Short:   "Run automation scripts",
		Aliases: []string{"run"},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "bash [script-path]",
		Short: "Run a bash script",
		Args:  cobra.ExactArgs(1),
		Run:   runBashScript,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list [directory]",
		Short: "List available scripts in directory",
		Args:  cobra.MaximumNArgs(1),
		Run:   listScripts,
	})

	return cmd
}

func runBashScript(cmd *cobra.Command, args []string) {
	scriptPath := args[0]

	if !fileExists(scriptPath) {
		log.Fatalf("Script not found: %s", scriptPath)
	}

	fmt.Printf("Running script: %s\n", scriptPath)
	bashCmd := exec.Command("bash", scriptPath)
	bashCmd.Stdout = os.Stdout
	bashCmd.Stderr = os.Stderr

	if err := bashCmd.Run(); err != nil {
		log.Fatalf("Script execution failed: %v", err)
	}
}

func listScripts(cmd *cobra.Command, args []string) {
	dir := getCurrentDir(args)

	files, err := filepath.Glob(filepath.Join(dir, "*.sh"))
	if err != nil {
		log.Fatalf("Error listing scripts: %v", err)
	}

	if len(files) == 0 {
		fmt.Printf("No bash scripts found in %s\n", dir)
		return
	}

	fmt.Printf("Available scripts in %s:\n", dir)
	for _, file := range files {
		fmt.Printf("  - %s\n", filepath.Base(file))
	}
}

// Pipeline commands
func pipelineCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pipeline",
		Short:   "Manage CI/CD pipelines",
		Aliases: []string{"ci"},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "create [template-name]",
		Short: "Create a CI/CD pipeline from template",
		Args:  cobra.ExactArgs(1),
		Run:   createPipeline,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "templates",
		Short: "List available pipeline templates",
		Run:   listPipelineTemplates,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "push",
		Short: "Push pipeline to repository",
		Run:   pushPipeline,
	})

	return cmd
}

func createPipeline(cmd *cobra.Command, args []string) {
	templateName := args[0]

	templates := getPipelineTemplates()
	template, exists := templates[templateName]
	if !exists {
		fmt.Printf("Template '%s' not found. Available templates:\n", templateName)
		for name := range templates {
			fmt.Printf("  - %s\n", name)
		}
		return
	}

	// Create pipeline directory
	pipelineDir := ".github/workflows"
	if err := os.MkdirAll(pipelineDir, 0755); err != nil {
		log.Fatalf("Failed to create pipeline directory: %v", err)
	}

	// Generate pipeline file
	pipelineFile := filepath.Join(pipelineDir, fmt.Sprintf("%s.yml", templateName))
	if err := generatePipelineFile(pipelineFile, template); err != nil {
		log.Fatalf("Failed to generate pipeline: %v", err)
	}

	fmt.Printf("Pipeline created: %s\n", pipelineFile)
}

func listPipelineTemplates(cmd *cobra.Command, args []string) {
	templates := getPipelineTemplates()

	fmt.Println("Available pipeline templates:")
	for name, template := range templates {
		fmt.Printf("  - %s (Language: %s)\n", name, template.Language)
	}
}

func pushPipeline(cmd *cobra.Command, args []string) {
	fmt.Print("Commit message: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	commitMsg := scanner.Text()

	if commitMsg == "" {
		commitMsg = "Add CI/CD pipeline"
	}

	// Git commands
	commands := [][]string{
		{"git", "add", ".github/workflows/"},
		{"git", "commit", "-m", commitMsg},
		{"git", "push", "origin", config.GitBranch},
	}

	for _, cmdArgs := range commands {
		gitCmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
		gitCmd.Stdout = os.Stdout
		gitCmd.Stderr = os.Stderr

		fmt.Printf("Running: %s\n", strings.Join(cmdArgs, " "))
		if err := gitCmd.Run(); err != nil {
			log.Fatalf("Git command failed: %v", err)
		}
	}

	fmt.Println("Pipeline pushed to repository successfully!")
}

// Configuration commands
func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "show",
		Short: "Show current configuration",
		Run:   showConfig,
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set configuration value",
		Args:  cobra.ExactArgs(2),
		Run:   setConfig,
	})

	return cmd
}

func showConfig(cmd *cobra.Command, args []string) {
	fmt.Println("Current configuration:")
	fmt.Printf("  Terraform Path: %s\n", config.TerraformPath)
	fmt.Printf("  Workspace Dir:  %s\n", config.WorkspaceDir)
	fmt.Printf("  Git Remote:     %s\n", config.GitRemote)
	fmt.Printf("  Git Branch:     %s\n", config.GitBranch)
}

func setConfig(cmd *cobra.Command, args []string) {
	key, value := args[0], args[1]

	switch key {
	case "terraform-path":
		config.TerraformPath = value
	case "workspace-dir":
		config.WorkspaceDir = value
	case "git-remote":
		config.GitRemote = value
	case "git-branch":
		config.GitBranch = value
	default:
		fmt.Printf("Unknown configuration key: %s\n", key)
		return
	}

	fmt.Printf("Set %s = %s\n", key, value)
}

// Helper functions
func getCurrentDir(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	pwd, _ := os.Getwd()
	return pwd
}

func confirmAction() bool {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))
	return response == "yes" || response == "y"
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func getPipelineTemplates() map[string]PipelineTemplate {
	return map[string]PipelineTemplate{
		"node": {
			Name:     "Node.js",
			Language: "javascript",
			BuildSteps: []string{
				"npm ci",
				"npm run build",
			},
			TestSteps: []string{
				"npm test",
				"npm run lint",
			},
			DeploySteps: []string{
				"npm run deploy",
			},
		},
		"go": {
			Name:     "Go",
			Language: "go",
			BuildSteps: []string{
				"go mod download",
				"go build -v ./...",
			},
			TestSteps: []string{
				"go test -v ./...",
				"go vet ./...",
			},
			DeploySteps: []string{
				"go build -o app",
			},
		},
		"python": {
			Name:     "Python",
			Language: "python",
			BuildSteps: []string{
				"pip install -r requirements.txt",
			},
			TestSteps: []string{
				"pytest",
				"flake8 .",
			},
			DeploySteps: []string{
				"python setup.py sdist bdist_wheel",
			},
		},
	}
}

func generatePipelineFile(filename string, template PipelineTemplate) error {
	githubActionsTemplate := `name: {{.Name}} CI/CD

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  build-and-test:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup {{.Language}}
      {{if eq .Language "javascript"}}uses: actions/setup-node@v3
      with:
        node-version: '18'{{end}}{{if eq .Language "go"}}uses: actions/setup-go@v3
      with:
        go-version: '1.19'{{end}}{{if eq .Language "python"}}uses: actions/setup-python@v3
      with:
        python-version: '3.9'{{end}}
    
    - name: Build
      run: |
        {{range .BuildSteps}}{{.}}
        {{end}}
    
    - name: Test
      run: |
        {{range .TestSteps}}{{.}}
        {{end}}
    
    - name: Deploy
      if: github.ref == 'refs/heads/main'
      run: |
        {{range .DeploySteps}}{{.}}
        {{end}}
`

	tmpl, err := tmpl.New("pipeline").Parse(githubActionsTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return tmpl.Execute(file, template)
}
