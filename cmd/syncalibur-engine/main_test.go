package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateShannonEntropy(t *testing.T) {
	// Normal low-entropy text
	normalText := "hello world this is a normal test string"
	entNormal := calculateShannonEntropy(normalText)
	if entNormal > 4.2 {
		t.Errorf("expected normal text entropy to be < 4.2, got %f", entNormal)
	}

	// High entropy random hash / secret
	randomSecret := "aZ9#kL2$vM8@xP4%qR7^wT1&yU5*"
	entSecret := calculateShannonEntropy(randomSecret)
	if entSecret < 4.4 {
		t.Errorf("expected random secret entropy to be >= 4.4, got %f", entSecret)
	}
}

func TestDetectSecretsInFile(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Clean file
	cleanFile := filepath.Join(tmpDir, "clean.go")
	_ = os.WriteFile(cleanFile, []byte("package main\n\nfunc main() {\n\tprintln(\"clean\")\n}\n"), 0644)
	if hasSecret, reason := detectSecretsInFile(cleanFile); hasSecret {
		t.Errorf("clean file flagged as secret: %s", reason)
	}

	// 2. File with AWS Key
	awsFile := filepath.Join(tmpDir, "aws.txt")
	_ = os.WriteFile(awsFile, []byte("AWS_KEY=AKIA1234567890ABCDEF\n"), 0644)
	if hasSecret, _ := detectSecretsInFile(awsFile); !hasSecret {
		t.Errorf("expected AWS key to be detected")
	}

	// 3. File with GitHub PAT
	ghFile := filepath.Join(tmpDir, "gh.txt")
	_ = os.WriteFile(ghFile, []byte("TOKEN=ghp_123456789012345678901234567890123456\n"), 0644)
	if hasSecret, _ := detectSecretsInFile(ghFile); !hasSecret {
		t.Errorf("expected GitHub PAT to be detected")
	}

	// 4. File with Private Key Header
	keyFile := filepath.Join(tmpDir, "priv.key")
	_ = os.WriteFile(keyFile, []byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----\n"), 0644)
	if hasSecret, _ := detectSecretsInFile(keyFile); !hasSecret {
		t.Errorf("expected Private Key header to be detected")
	}

	// 5. Benign File with UUID (should NOT trigger secret quarantine)
	uuidFile := filepath.Join(tmpDir, "uuid.txt")
	_ = os.WriteFile(uuidFile, []byte("SESSION_ID=c82c1970-4ff7-461a-80d0-e958b320a7ec\n"), 0644)
	if hasSecret, reason := detectSecretsInFile(uuidFile); hasSecret {
		t.Errorf("expected UUID to be ignored, got false alarm: %s", reason)
	}

	// 6. Benign File with base64 data URI (should NOT trigger secret quarantine)
	dataUriFile := filepath.Join(tmpDir, "asset.js")
	_ = os.WriteFile(dataUriFile, []byte("const icon = \"data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==\";\n"), 0644)
	if hasSecret, reason := detectSecretsInFile(dataUriFile); hasSecret {
		t.Errorf("expected data URI to be ignored, got false alarm: %s", reason)
	}
}

func TestValidateFileSyntax(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Valid Go file
	validGo := filepath.Join(tmpDir, "valid.go")
	_ = os.WriteFile(validGo, []byte("package test\nfunc Add(a, b int) int { return a + b }\n"), 0644)
	if valid, errReason := validateFileSyntax(tmpDir, "valid.go"); !valid {
		t.Errorf("expected valid Go file to pass syntax check, got: %s", errReason)
	}

	// 2. Broken Go file (unclosed brace)
	brokenGo := filepath.Join(tmpDir, "broken.go")
	_ = os.WriteFile(brokenGo, []byte("package test\nfunc Broken() {\n"), 0644)
	if valid, _ := validateFileSyntax(tmpDir, "broken.go"); valid {
		t.Errorf("expected broken Go file (missing brace) to fail syntax check")
	}

	// 3. Valid JSON
	validJSON := filepath.Join(tmpDir, "valid.json")
	_ = os.WriteFile(validJSON, []byte(`{"name": "syncalibur", "active": true}`), 0644)
	if valid, errReason := validateFileSyntax(tmpDir, "valid.json"); !valid {
		t.Errorf("expected valid JSON to pass, got: %s", errReason)
	}

	// 4. Broken JSON
	brokenJSON := filepath.Join(tmpDir, "broken.json")
	_ = os.WriteFile(brokenJSON, []byte(`{"name": "syncalibur", broken`), 0644)
	if valid, _ := validateFileSyntax(tmpDir, "broken.json"); valid {
		t.Errorf("expected broken JSON to fail syntax check")
	}

	// 5. Valid Bash script with case statement and hash comments
	validSh := filepath.Join(tmpDir, "valid.sh")
	_ = os.WriteFile(validSh, []byte("#!/bin/bash\n# comment with (unclosed paren\ncase \"$1\" in\n  start) echo ok ;;\n  stop)  echo bye ;;\nesac\n"), 0755)
	if valid, errReason := validateFileSyntax(tmpDir, "valid.sh"); !valid {
		t.Errorf("expected valid Bash script with case and comment to pass, got: %s", errReason)
	}

	// 6. Broken Bash script (unclosed if statement)
	brokenSh := filepath.Join(tmpDir, "broken.sh")
	_ = os.WriteFile(brokenSh, []byte("#!/bin/bash\nif [ 1 -eq 1 ]; then echo ok\n"), 0755)
	if valid, _ := validateFileSyntax(tmpDir, "broken.sh"); valid {
		t.Errorf("expected broken Bash script to fail syntax check")
	}

	// 7. Valid Java file
	validJava := filepath.Join(tmpDir, "App.java")
	_ = os.WriteFile(validJava, []byte("public class App { public static void main(String[] args) { System.out.println(\"ok\"); } }"), 0644)
	if valid, errReason := validateFileSyntax(tmpDir, "App.java"); !valid {
		t.Errorf("expected valid Java file to pass, got: %s", errReason)
	}

	// 8. Broken Java file (missing closing brace)
	brokenJava := filepath.Join(tmpDir, "BrokenApp.java")
	_ = os.WriteFile(brokenJava, []byte("public class BrokenApp { public static void main(String[] args) {"), 0644)
	if valid, _ := validateFileSyntax(tmpDir, "BrokenApp.java"); valid {
		t.Errorf("expected broken Java file to fail syntax check")
	}

	// 9. Valid React TSX
	validTSX := filepath.Join(tmpDir, "Component.tsx")
	_ = os.WriteFile(validTSX, []byte("export const Header = () => { return <div><h1>Title</h1></div>; };"), 0644)
	if valid, errReason := validateFileSyntax(tmpDir, "Component.tsx"); !valid {
		t.Errorf("expected valid React TSX to pass, got: %s", errReason)
	}

	// 10. Broken React TSX (unclosed JSX expression brace)
	brokenTSX := filepath.Join(tmpDir, "BrokenComponent.tsx")
	_ = os.WriteFile(brokenTSX, []byte("export const Header = () => { return <div>{items.map(i => <span>{i}</span>); };"), 0644)
	if valid, _ := validateFileSyntax(tmpDir, "BrokenComponent.tsx"); valid {
		t.Errorf("expected broken React TSX to fail syntax check")
	}

	// 11. Valid Python with hash comment and multiline docstring
	validPy := filepath.Join(tmpDir, "script.py")
	pyContent := `# Python script with (unmatched parenthesis in comment
"""
Docstring with [unmatched brackets] and (parentheses)
"""
def hello():
    return "world"
`
	_ = os.WriteFile(validPy, []byte(pyContent), 0644)
	if valid, errReason := validateFileSyntax(tmpDir, "script.py"); !valid {
		t.Errorf("expected valid Python with docstring and comments to pass, got: %s", errReason)
	}

	// 12. Broken Python (unclosed bracket in code)
	brokenPy := filepath.Join(tmpDir, "broken.py")
	_ = os.WriteFile(brokenPy, []byte("def broken():\n    return [1, 2, 3\n"), 0644)
	if valid, _ := validateFileSyntax(tmpDir, "broken.py"); valid {
		t.Errorf("expected broken Python to fail syntax check")
	}
}

func TestInspectGoAST(t *testing.T) {
	tmpDir := t.TempDir()
	sampleGo := filepath.Join(tmpDir, "sample.go")
	src := `package sample

type AuthEngine struct{}

func (a *AuthEngine) ValidateToken(t string) bool {
	return t != ""
}

func HelperFunc() int {
	return 42
}
`
	_ = os.WriteFile(sampleGo, []byte(src), 0644)

	// Check line 6 (inside ValidateToken method)
	res := inspectGoAST(sampleGo, 6, 6)
	if res == nil {
		t.Fatalf("expected AST inspection to find symbol at line 6")
	}
	if res.Kind != "method" || res.Symbol != "(*AuthEngine).ValidateToken" {
		t.Errorf("expected (*AuthEngine).ValidateToken, got kind=%s, symbol=%s", res.Kind, res.Symbol)
	}

	// Check line 10 (inside HelperFunc)
	resFunc := inspectGoAST(sampleGo, 10, 10)
	if resFunc == nil {
		t.Fatalf("expected AST inspection to find symbol at line 10")
	}
	if resFunc.Kind != "func" || resFunc.Symbol != "HelperFunc()" {
		t.Errorf("expected HelperFunc(), got kind=%s, symbol=%s", resFunc.Kind, resFunc.Symbol)
	}
}
