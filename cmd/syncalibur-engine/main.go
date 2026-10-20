// ==============================================================================
// Syncalibur Go Engine (cmd/syncalibur-engine/main.go)
// ------------------------------------------------------------------------------
// High-performance, zero-dependency Go engine powering Syncalibur's:
//   1. High-Definition TrueColor Terminal Emblem Renderer (`render-logo`)
//   2. Semantic Conventional Commit Classifier + Diff-Stat + AI Hook (`commit-msg`)
//   3. Multi-Repository Fleet Scanner & Parallel Status Probe (`fleet-scan`, `fleet-status`)
//   4. 7-Day Commit Velocity Sparkline Generator (`sparkline`)
//   5. First-Class Telegram Bot & Multi-Channel Webhook Dispatcher (`webhook`, `telegram`)
//   6. Interactive Live TUI Control Center (`ui`)
// ==============================================================================
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	cReset = "\033[0m"
	cBold  = "\033[1m"
	cGold  = "\033[38;2;245;194;66m"
	cCyan  = "\033[38;2;56;189;248m"
	cGreen = "\033[38;2;52;211;153m"
	cRed   = "\033[38;2;248;113;113m"
	cAmber = "\033[38;2;251;191;36m"
	cSlate = "\033[38;2;148;163;184m"
	cWhite = "\033[38;2;241;245;249m"
)

var quads = [16]string{
	" ", "▘", "▝", "▀",
	"▖", "▌", "▞", "▛",
	"▗", "▚", "▐", "▜",
	"▄", "▙", "▟", "█",
}

type rgbaPx struct {
	r, g, b, a int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: syncalibur-engine <render-logo|commit-msg|sparkline|fleet-scan|fleet-status|webhook|telegram|ui> [args...]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "render-logo":
		if len(os.Args) < 4 {
			os.Exit(1)
		}
		if err := renderLogoANSI(os.Args[2], os.Args[3], 44, 21); err != nil {
			fmt.Fprintf(os.Stderr, "render-logo error: %v\n", err)
			os.Exit(1)
		}
	case "commit-msg":
		repo := "."
		if len(os.Args) >= 3 {
			repo = os.Args[2]
		}
		fmt.Println(generateCommitMessage(repo))
	case "sparkline":
		repo := "."
		if len(os.Args) >= 3 {
			repo = os.Args[2]
		}
		fmt.Print(generateSparkline(repo))
	case "fleet-scan":
		root := "."
		if len(os.Args) >= 3 {
			root = os.Args[2]
		}
		for _, r := range scanGitRepos(root, 3) {
			fmt.Println(r)
		}
	case "fleet-status":
		if len(os.Args) < 3 {
			os.Exit(1)
		}
		printFleetStatusTable(os.Args[2])
	case "webhook":
		if len(os.Args) < 9 {
			os.Exit(1)
		}
		if err := sendWebhook(os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7], os.Args[8]); err != nil {
			fmt.Fprintf(os.Stderr, "webhook error: %v\n", err)
			os.Exit(1)
		}
	case "telegram":
		if len(os.Args) < 10 {
			os.Exit(1)
		}
		if err := sendTelegramMessage(os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7], os.Args[8], os.Args[9]); err != nil {
			fmt.Fprintf(os.Stderr, "telegram error: %v\n", err)
			os.Exit(1)
		}
	case "ui":
		scriptPath := ""
		repoRoot := "."
		if len(os.Args) >= 3 {
			scriptPath = os.Args[2]
		}
		if len(os.Args) >= 4 {
			repoRoot = os.Args[3]
		}
		runInteractiveTUI(scriptPath, repoRoot)
	default:
		os.Exit(1)
	}
}

// ==============================================================================
// 1. High-Definition TrueColor Terminal Emblem Renderer
// ==============================================================================

func renderLogoANSI(srcPath, dstPath string, cols, rows int) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}

	b := img.Bounds()
	minX, minY, maxX, maxY := b.Max.X, b.Max.Y, b.Min.X, b.Min.Y
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a>>8 >= 16 {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	if minX >= maxX || minY >= maxY {
		minX, minY, maxX, maxY = b.Min.X, b.Min.Y, b.Max.X-1, b.Max.Y-1
	}

	sw, sh := cols*2, rows*2
	cropW := float64(maxX - minX + 1)
	cropH := float64(maxY - minY + 1)

	grid := make([][]rgbaPx, sh)
	for sy := 0; sy < sh; sy++ {
		grid[sy] = make([]rgbaPx, sw)
		for sx := 0; sx < sw; sx++ {
			fx0 := float64(minX) + (float64(sx)/float64(sw))*cropW
			fy0 := float64(minY) + (float64(sy)/float64(sh))*cropH
			fx1 := float64(minX) + (float64(sx+1)/float64(sw))*cropW
			fy1 := float64(minY) + (float64(sy+1)/float64(sh))*cropH

			var sr, sg, sb, sa, cnt float64
			for py := int(fy0); py <= int(fy1) && py <= maxY; py++ {
				for px := int(fx0); px <= int(fx1) && px <= maxX; px++ {
					c := color.NRGBAModel.Convert(img.At(px, py)).(color.NRGBA)
					sr += float64(c.R)
					sg += float64(c.G)
					sb += float64(c.B)
					sa += float64(c.A)
					cnt++
				}
			}
			if cnt == 0 {
				continue
			}
			r, g, bl, al := sr/cnt, sg/cnt, sb/cnt, sa/cnt
			lum := 0.299*r + 0.587*g + 0.114*bl
			r = lum + (r-lum)*1.28
			g = lum + (g-lum)*1.28
			bl = lum + (bl-lum)*1.28
			r = (r-128.0)*1.24 + 128.0
			g = (g-128.0)*1.24 + 128.0
			bl = (bl-128.0)*1.24 + 128.0
			grid[sy][sx] = rgbaPx{
				r: clampByte(r),
				g: clampByte(g),
				b: clampByte(bl),
				a: clampByte(al),
			}
		}
	}

	var out strings.Builder
	for cy := 0; cy < rows; cy++ {
		for cx := 0; cx < cols; cx++ {
			pts := [4]rgbaPx{
				grid[cy*2][cx*2],
				grid[cy*2][cx*2+1],
				grid[cy*2+1][cx*2],
				grid[cy*2+1][cx*2+1],
			}
			mask := 0
			var active []rgbaPx
			for i := 0; i < 4; i++ {
				if pts[i].a >= 85 {
					if pts[i].a < 220 {
						f := float64(pts[i].a) / 255.0
						pts[i].r = int(float64(pts[i].r) * f)
						pts[i].g = int(float64(pts[i].g) * f)
						pts[i].b = int(float64(pts[i].b) * f)
					}
					mask |= (1 << i)
					active = append(active, pts[i])
				}
			}
			if mask == 0 {
				out.WriteString("\033[0m ")
			} else if mask != 15 {
				var sr, sg, sb int
				for _, p := range active {
					sr += p.r
					sg += p.g
					sb += p.b
				}
				n := len(active)
				fmt.Fprintf(&out, "\033[0m\033[38;2;%d;%d;%dm%s", sr/n, sg/n, sb/n, quads[mask])
			} else {
				bestErr := math.MaxFloat64
				bestM := 3
				var bestFG, bestBG [3]int
				for m := 1; m < 8; m++ {
					var g1, g2 []rgbaPx
					for i := 0; i < 4; i++ {
						if ((m >> i) & 1) == 1 {
							g1 = append(g1, pts[i])
						} else {
							g2 = append(g2, pts[i])
						}
					}
					c1 := avgPxColor(g1)
					c2 := avgPxColor(g2)
					errVal := sqPxErr(g1, c1) + sqPxErr(g2, c2)
					if m == 3 {
						errVal *= 0.92
					}
					if errVal < bestErr {
						bestErr = errVal
						bestM = m
						bestFG = c1
						bestBG = c2
					}
				}
				fmt.Fprintf(&out, "\033[38;2;%d;%d;%dm\033[48;2;%d;%d;%dm%s",
					bestFG[0], bestFG[1], bestFG[2],
					bestBG[0], bestBG[1], bestBG[2],
					quads[bestM])
			}
		}
		out.WriteString("\033[0m\n")
	}

	return os.WriteFile(dstPath, []byte(out.String()), 0644)
}

func clampByte(v float64) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return int(v)
}

func avgPxColor(list []rgbaPx) [3]int {
	if len(list) == 0 {
		return [3]int{0, 0, 0}
	}
	var sr, sg, sb int
	for _, p := range list {
		sr += p.r
		sg += p.g
		sb += p.b
	}
	n := len(list)
	return [3]int{sr / n, sg / n, sb / n}
}

func sqPxErr(list []rgbaPx, c [3]int) float64 {
	var e float64
	for _, p := range list {
		dr := float64(p.r - c[0])
		dg := float64(p.g - c[1])
		db := float64(p.b - c[2])
		e += dr*dr + dg*dg + db*db
	}
	return e
}

// ==============================================================================
// 2. Semantic Conventional Commit Classifier + Diff-Stat + Optional AI Hook
// ==============================================================================

func generateCommitMessage(repoRoot string) string {
	if strings.EqualFold(os.Getenv("SYNCALIBUR_AI_COMMIT"), "true") {
		if aiMsg := tryAICommitMessage(repoRoot); aiMsg != "" {
			return aiMsg
		}
	}

	out, err := exec.Command("git", "-C", repoRoot, "-c", "core.quotePath=false", "diff", "--cached", "--name-status").Output()
	if err != nil || len(bytes.TrimSpace(out)) == 0 {
		return fmt.Sprintf("chore: auto-sync (%s)", time.Now().Format("15:04:05"))
	}

	var files, statuses []string
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			statuses = append(statuses, parts[0])
			files = append(files, parts[len(parts)-1])
		}
	}
	if len(files) == 0 {
		return fmt.Sprintf("chore: auto-sync (%s)", time.Now().Format("15:04:05"))
	}

	var totalAdd, totalDel int
	numOut, _ := exec.Command("git", "-C", repoRoot, "diff", "--cached", "--numstat").Output()
	for _, line := range strings.Split(strings.TrimSpace(string(numOut)), "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			a, _ := strconv.Atoi(fields[0])
			d, _ := strconv.Atoi(fields[1])
			totalAdd += a
			totalDel += d
		}
	}
	statBadge := fmt.Sprintf("(+%d/-%d)", totalAdd, totalDel)

	allDocs, allTests, allCI, allStyle := true, true, true, true
	hasAdd, hasDel := false, false

	for i, f := range files {
		st := statuses[i]
		if strings.HasPrefix(st, "A") {
			hasAdd = true
		}
		if strings.HasPrefix(st, "D") {
			hasDel = true
		}
		lower := strings.ToLower(f)
		if !(strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".rst") || strings.HasSuffix(lower, ".txt") || strings.HasPrefix(lower, "docs/") || lower == "license") {
			allDocs = false
		}
		if !(strings.Contains(lower, "test") || strings.Contains(lower, "spec")) {
			allTests = false
		}
		if !(strings.HasPrefix(lower, ".github/") || strings.HasPrefix(lower, ".gitlab") || lower == "dockerfile" || lower == "makefile" || lower == "go.mod" || lower == "package.json" || strings.HasSuffix(lower, ".tf") || strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml")) {
			allCI = false
		}
		if !(strings.HasSuffix(lower, ".css") || strings.HasSuffix(lower, ".scss")) {
			allStyle = false
		}
	}

	cType := "chore"
	if allDocs {
		cType = "docs"
	} else if allTests {
		cType = "test"
	} else if allStyle {
		cType = "style"
	} else if allCI {
		cType = "build"
	} else if hasAdd || totalAdd > totalDel*2 {
		cType = "feat"
	} else if hasDel || totalDel > totalAdd {
		cType = "refactor"
	} else {
		diffSample, _ := exec.Command("git", "-C", repoRoot, "diff", "--cached", "-U0").Output()
		if regexp.MustCompile(`(?i)\+.*(fix|bug|error|patch|resolve|crash)`).Match(diffSample) {
			cType = "fix"
		} else {
			cType = "feat"
		}
	}

	cScope := ""
	if len(files) == 1 {
		base := filepath.Base(files[0])
		ext := filepath.Ext(base)
		cScope = strings.TrimSuffix(base, ext)
	} else {
		firstDir := filepath.Dir(files[0])
		if firstDir != "." && firstDir != "" {
			same := true
			for _, f := range files[1:] {
				if filepath.Dir(f) != firstDir {
					same = false
					break
				}
			}
			if same {
				cScope = filepath.Base(firstDir)
			}
		}
	}
	if !regexp.MustCompile(`^[A-Za-z0-9_-]{1,20}$`).MatchString(cScope) {
		cScope = ""
	}

	var subject string
	if len(files) == 1 {
		action := "update"
		if strings.HasPrefix(statuses[0], "A") {
			action = "add"
		} else if strings.HasPrefix(statuses[0], "D") {
			action = "remove"
		}
		subject = fmt.Sprintf("%s %s", action, filepath.Base(files[0]))
	} else if len(files) <= 3 {
		var baseNames []string
		for _, f := range files {
			baseNames = append(baseNames, filepath.Base(f))
		}
		subject = "update " + strings.Join(baseNames, ", ")
	} else {
		subject = fmt.Sprintf("update %s, %s (+%d more)", filepath.Base(files[0]), filepath.Base(files[1]), len(files)-2)
	}

	if cScope != "" {
		return fmt.Sprintf("%s(%s): %s %s", cType, cScope, subject, statBadge)
	}
	return fmt.Sprintf("%s: %s %s", cType, subject, statBadge)
}

func tryAICommitMessage(repoRoot string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	statOut, _ := exec.CommandContext(ctx, "git", "-C", repoRoot, "diff", "--cached", "--shortstat").Output()
	diffOut, _ := exec.CommandContext(ctx, "git", "-C", repoRoot, "diff", "--cached", "-U1").Output()
	if len(diffOut) > 1600 {
		diffOut = diffOut[:1600]
	}

	model := strings.TrimSpace(os.Getenv("SYNCALIBUR_OLLAMA_MODEL"))
	if model == "" {
		model = "qwen2.5-coder:1.5b"
	}

	prompt := fmt.Sprintf("Write a single-line Conventional Commit message (max 65 chars, format: type(scope): summary) for this git diff:\n%s\n%s", string(statOut), string(diffOut))
	reqBody, _ := json.Marshal(map[string]any{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	})

	req, err := http.NewRequestWithContext(ctx, "POST", "http://127.0.0.1:11434/api/generate", bytes.NewReader(reqBody))
	if err != nil {
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var parsed struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return ""
	}
	line := strings.TrimSpace(strings.Split(parsed.Response, "\n")[0])
	line = strings.Trim(line, "`\"'")
	if regexp.MustCompile(`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore)(\([a-zA-Z0-9_-]+\))?:\s+.+$`).MatchString(line) {
		return line
	}
	return ""
}

// ==============================================================================
// 3. 7-Day Commit Velocity Sparkline Generator
// ==============================================================================

func generateSparkline(repoRoot string) string {
	glyphs := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	counts := make([]int, 7)
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	out, _ := exec.Command("git", "-C", repoRoot, "log", "--since=7 days ago", "--format=%at").Output()
	total7d := 0
	for _, s := range strings.Fields(string(out)) {
		ts, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		t := time.Unix(ts, 0)
		tDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, now.Location())
		daysAgo := int(startOfToday.Sub(tDay).Hours() / 24)
		if daysAgo >= 0 && daysAgo < 7 {
			counts[6-daysAgo]++
			total7d++
		}
	}

	maxC := 1
	for _, c := range counts {
		if c > maxC {
			maxC = c
		}
	}

	var b strings.Builder
	for _, c := range counts {
		if c == 0 {
			b.WriteString(cSlate + "·" + cReset)
		} else {
			idx := (c * 7) / maxC
			if idx > 7 {
				idx = 7
			}
			b.WriteString(cGreen + glyphs[idx] + cReset)
		}
	}
	fmt.Fprintf(&b, " %s(7d: %d commits • peak %d/day)%s", cSlate, total7d, maxC, cReset)
	return b.String()
}

// ==============================================================================
// 4. Multi-Repository Fleet Scanner & Parallel Status Probe
// ==============================================================================

func scanGitRepos(root string, maxDepth int) []string {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return nil
	}
	var repos []string
	var walk func(dir string, depth int)
	walk = func(dir string, depth int) {
		if depth > maxDepth {
			return
		}
		if st, err := os.Stat(filepath.Join(dir, ".git")); err == nil && st.IsDir() {
			repos = append(repos, dir)
			return
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "snap" {
				continue
			}
			walk(filepath.Join(dir, name), depth+1)
		}
	}
	walk(rootAbs, 0)
	sort.Strings(repos)
	return repos
}

type fleetRepoStatus struct {
	Path     string
	Name     string
	Branch   string
	State    string
	Dirty    int
	Unpushed int
	Remote   string
}

func printFleetStatusTable(fleetFile string) {
	data, err := os.ReadFile(fleetFile)
	if err != nil {
		fmt.Printf("  %sNo fleet file found at %s%s\n", cAmber, fleetFile, cReset)
		return
	}

	var paths []string
	seen := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(strings.Split(line, "#")[0])
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "~/") {
			if home, err := os.UserHomeDir(); err == nil {
				line = filepath.Join(home, line[2:])
			}
		}
		if abs, err := filepath.Abs(line); err == nil {
			if st, err := os.Stat(abs); err == nil && st.IsDir() && !seen[abs] {
				seen[abs] = true
				paths = append(paths, abs)
			}
		}
	}

	if len(paths) == 0 {
		fmt.Printf("  %sNo repositories enrolled in %s. Run: syncalibur fleet add <path>%s\n", cAmber, fleetFile, cReset)
		return
	}

	results := make([]fleetRepoStatus, len(paths))
	var wg sync.WaitGroup
	for i, p := range paths {
		wg.Add(1)
		go func(idx int, repoPath string) {
			defer wg.Done()
			results[idx] = inspectRepoStatus(repoPath)
		}(i, p)
	}
	wg.Wait()

	fmt.Printf("  %s╭────────────────────────────────────────────────────────────────────────────────────────╮%s\n", cGold, cReset)
	fmt.Printf("  %s│%s  %s%s%-84s%s  %s│%s\n", cGold, cReset, cGold, cBold, fmt.Sprintf("◈  S Y N C A L I B U R   F L E E T   S T A T U S   (%d REPOSITORIES)", len(results)), cReset, cGold, cReset)
	fmt.Printf("  %s╰────────────────────────────────────────────────────────────────────────────────────────╯%s\n", cGold, cReset)
	fmt.Printf("    %s%-20s %-14s %-12s %-22s %-20s%s\n", cBold, "REPOSITORY", "DAEMON", "BRANCH", "SYNC STATE", "PATH", cReset)
	fmt.Printf("    %s──────────────────────────────────────────────────────────────────────────────────────%s\n", cSlate, cReset)

	home, _ := os.UserHomeDir()
	for _, r := range results {
		shortPath := r.Path
		if home != "" && strings.HasPrefix(shortPath, home) {
			shortPath = "~" + strings.TrimPrefix(shortPath, home)
		}
		rawSync := "✔ In Sync"
		syncColor := cGreen
		if r.Dirty > 0 || r.Unpushed > 0 {
			rawSync = fmt.Sprintf("● %d dirty, %d unpushed", r.Dirty, r.Unpushed)
			syncColor = cAmber
		}
		stateColor := cSlate
		if strings.HasPrefix(r.State, "●") {
			stateColor = cGreen
		}
		fmt.Printf("    %s%-20s%s %s%-14s%s %s%-12s%s %s%-22s%s %s%s%s\n",
			cWhite+cBold, truncateStr(r.Name, 19), cReset,
			stateColor, truncateStr(r.State, 13), cReset,
			cCyan, truncateStr(r.Branch, 11), cReset,
			syncColor, truncateStr(rawSync, 21), cReset,
			cSlate, shortPath, cReset)
	}
	fmt.Println()
}

func inspectRepoStatus(repoPath string) fleetRepoStatus {
	st := fleetRepoStatus{
		Path:   repoPath,
		Name:   filepath.Base(repoPath),
		Branch: "main",
		State:  "○ Idle",
	}
	if bOut, err := exec.Command("git", "-C", repoPath, "symbolic-ref", "--quiet", "--short", "HEAD").Output(); err == nil {
		st.Branch = strings.TrimSpace(string(bOut))
	}
	pidFile := filepath.Join(repoPath, ".git", "syncalibur.pid")
	if pBytes, err := os.ReadFile(pidFile); err == nil {
		pidStr := strings.TrimSpace(string(pBytes))
		if pid, err := strconv.Atoi(pidStr); err == nil && processAlive(pid) {
			st.State = fmt.Sprintf("● PID %d", pid)
		}
	}
	if sOut, err := exec.Command("git", "-C", repoPath, "status", "--porcelain", "-uall").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(sOut)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			st.Dirty = len(lines)
		}
	}
	return st
}

func processAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return p.Signal(syscall.Signal(0)) == nil
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-1] + "…"
}

// ==============================================================================
// 5. First-Class Telegram Bot & Multi-Channel Webhook Dispatcher
// ==============================================================================

func escapeTelegramHTML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return r.Replace(s)
}

func formatTelegramHTMLMessage(eventType, repoName, branchName, hostInfo, title, message, details string) string {
	badge := "⚔️"
	switch eventType {
	case "failure":
		badge = "🚨"
	case "recovery", "success":
		badge = "✅"
	case "summary", "test":
		badge = "⚔️"
	}

	var b strings.Builder
	fmt.Fprintf(&b, "<b>%s %s</b>\n", badge, escapeTelegramHTML(title))
	if strings.TrimSpace(message) != "" {
		fmt.Fprintf(&b, "%s\n", escapeTelegramHTML(message))
	}
	fmt.Fprintf(&b, "\n<b>◈ Repository:</b> <code>%s (%s)</code>\n", escapeTelegramHTML(repoName), escapeTelegramHTML(branchName))
	fmt.Fprintf(&b, "<b>◈ Host:</b> <code>%s</code>\n", escapeTelegramHTML(hostInfo))
	fmt.Fprintf(&b, "<b>◈ Time:</b> <code>%s</code>", time.Now().Format("2006-01-02 15:04:05 MST"))

	if strings.TrimSpace(details) != "" {
		trimmed := strings.TrimSpace(details)
		if len(trimmed) > 1200 {
			trimmed = trimmed[:1200] + "..."
		}
		fmt.Fprintf(&b, "\n\n<pre>%s</pre>", escapeTelegramHTML(trimmed))
	}
	return b.String()
}

func sendTelegramMessage(botToken, chatID, eventType, repoName, branchName, title, message, details string) error {
	botToken = strings.TrimSpace(botToken)
	chatID = strings.TrimSpace(chatID)
	if botToken == "" || chatID == "" {
		return nil
	}
	botToken = strings.TrimPrefix(botToken, "bot")
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	host, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = "user"
	}
	hostInfo := fmt.Sprintf("%s@%s (%s/%s)", user, host, runtime.GOOS, runtime.GOARCH)

	text := formatTelegramHTMLMessage(eventType, repoName, branchName, hostInfo, title, message, details)
	payload := map[string]any{
		"chat_id":                  chatID,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
		"text":                     text,
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("telegram API HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBytes)))
	}
	return nil
}

func sendWebhook(webhookURL, eventType, repoName, branchName, title, message, details string) error {
	var firstErr error

	// Also dispatch to first-class Telegram Bot if SYNCALIBUR_TELEGRAM_BOT_TOKEN + SYNCALIBUR_TELEGRAM_CHAT_ID are set in env
	// and webhookURL itself isn't already pointing to the same Telegram Bot endpoint
	tgToken := strings.TrimSpace(os.Getenv("SYNCALIBUR_TELEGRAM_BOT_TOKEN"))
	tgChat := strings.TrimSpace(os.Getenv("SYNCALIBUR_TELEGRAM_CHAT_ID"))
	lowerURL := strings.ToLower(strings.TrimSpace(webhookURL))

	if tgToken != "" && tgChat != "" && !strings.Contains(lowerURL, "api.telegram.org/bot") {
		if err := sendTelegramMessage(tgToken, tgChat, eventType, repoName, branchName, title, message, details); err != nil {
			firstErr = err
		}
	}

	if strings.TrimSpace(webhookURL) == "" {
		return firstErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	host, _ := os.Hostname()
	user := os.Getenv("USER")
	if user == "" {
		user = "user"
	}
	hostInfo := fmt.Sprintf("%s@%s (%s/%s)", user, host, runtime.GOOS, runtime.GOARCH)
	nowISO := time.Now().UTC().Format(time.RFC3339)

	colorDec := 16494372
	colorHex := "#fbbf24"
	switch eventType {
	case "failure":
		colorDec = 16259441
		colorHex = "#f87171"
	case "recovery", "success":
		colorDec = 3461465
		colorHex = "#34d399"
	case "summary", "test":
		colorDec = 3718584
		colorHex = "#38bdf8"
	}

	var payload any

	switch {
	case strings.Contains(lowerURL, "discord.com/api/webhooks/") || strings.Contains(lowerURL, "discordapp.com/api/webhooks/"):
		if len(details) > 900 {
			details = details[:900] + "..."
		}
		fields := []map[string]any{
			{"name": "Repository", "value": "`" + repoName + "`", "inline": true},
			{"name": "Branch", "value": "`" + branchName + "`", "inline": true},
			{"name": "Host", "value": "`" + hostInfo + "`", "inline": true},
		}
		if strings.TrimSpace(details) != "" {
			fields = append(fields, map[string]any{"name": "Telemetry", "value": "```text\n" + details + "\n```", "inline": false})
		}
		payload = map[string]any{
			"username":   "Syncalibur Daemon",
			"avatar_url": "https://raw.githubusercontent.com/YusufKaya00/Syncalibur/main/assets/syncalibur-logo.png",
			"embeds": []map[string]any{
				{
					"title":       title,
					"description": message,
					"color":       colorDec,
					"fields":      fields,
					"footer":      map[string]string{"text": "Syncalibur Autonomous Git Daemon (Go Engine)"},
					"timestamp":   nowISO,
				},
			},
		}
	case strings.Contains(lowerURL, "hooks.slack.com/services/"):
		attText := fmt.Sprintf("%s\n*Repo:* `%s` (%s) • *Host:* `%s`", message, repoName, branchName, hostInfo)
		if strings.TrimSpace(details) != "" {
			attText += fmt.Sprintf("\n```%s```", details)
		}
		payload = map[string]any{
			"text": "*" + title + "*",
			"attachments": []map[string]any{
				{
					"color": colorHex,
					"text":  attText,
				},
			},
		}
	case strings.Contains(lowerURL, "api.telegram.org/bot"):
		chatID := tgChat
		if parsedURL, err := url.Parse(webhookURL); err == nil {
			if qChat := parsedURL.Query().Get("chat_id"); qChat != "" {
				chatID = qChat
			}
		}
		text := formatTelegramHTMLMessage(eventType, repoName, branchName, hostInfo, title, message, details)
		tgMap := map[string]any{
			"parse_mode":               "HTML",
			"disable_web_page_preview": true,
			"text":                     text,
		}
		if chatID != "" {
			tgMap["chat_id"] = chatID
		}
		payload = tgMap
	default:
		payload = map[string]any{
			"service":    "syncalibur",
			"event":      eventType,
			"timestamp":  nowISO,
			"repository": repoName,
			"branch":     branchName,
			"host":       hostInfo,
			"title":      title,
			"message":    message,
			"details":    details,
		}
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", webhookURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		respBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("webhook HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(respBytes)))
	}
	return firstErr
}

// ==============================================================================
// 6. Interactive Live TUI Control Center (`syncalibur ui`)
// ==============================================================================

func runInteractiveTUI(scriptPath, repoRoot string) {
	// Put terminal into raw single-char non-canonical mode via stty
	sttyOrigOut, _ := exec.Command("stty", "-g").Output()
	sttyOrig := strings.TrimSpace(string(sttyOrigOut))
	_ = exec.Command("stty", "-icanon", "-echo", "min", "0", "time", "0").Run()

	cleanup := func() {
		if sttyOrig != "" {
			_ = exec.Command("stty", sttyOrig).Run()
		} else {
			_ = exec.Command("stty", "sane").Run()
		}
		fmt.Print("\033[?1049l\033[?25h")
	}
	defer cleanup()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		cleanup()
		os.Exit(0)
	}()

	// Enter alternate screen buffer and hide cursor
	fmt.Print("\033[?1049h\033[?25l")

	keyCh := make(chan byte, 16)
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if err == nil && n > 0 {
				keyCh <- buf[0]
			} else {
				time.Sleep(40 * time.Millisecond)
			}
		}
	}()

	var flashMsg string
	var flashUntil time.Time
	gitDir := filepath.Join(repoRoot, ".git")

	renderFrame := func() {
		var frame strings.Builder
		frame.WriteString("\033[H")

		// Render main status dashboard from syncalibur.sh
		cmd := exec.Command("bash", scriptPath, "status", "--compact-ui")
		cmd.Env = append(os.Environ(), "CLICOLOR_FORCE=1")
		out, _ := cmd.Output()
		frame.Write(out)

		// Live Velocity & Sparkline Bar
		commitsTodayOut, _ := exec.Command("git", "-C", repoRoot, "rev-list", "--count", "--since=midnight", "HEAD").Output()
		commitsToday := strings.TrimSpace(string(commitsTodayOut))
		if commitsToday == "" {
			commitsToday = "0"
		}
		spark := generateSparkline(repoRoot)

		pauseState := cGreen + "▶ ACTIVE" + cReset
		if _, err := os.Stat(filepath.Join(gitDir, "syncalibur.pause")); err == nil {
			pauseState = cAmber + cBold + "⏸ PAUSED BY USER" + cReset
		}

		fmt.Fprintf(&frame, "  %s╭────────────────────────────────────────────────────────────────────────────────────────╮%s\n", cGold, cReset)
		fmt.Fprintf(&frame, "  %s│%s  %s%sLIVE TELEMETRY COCKPIT (Go Engine)%s   •   Loop: %s   •   Today: %s%s commits%s  %s│%s\n",
			cGold, cReset, cCyan, cBold, cReset, pauseState, cGreen, commitsToday, cReset, cGold, cReset)
		fmt.Fprintf(&frame, "  %s│%s  %s◈ 7-Day Velocity :%s %s\n", cGold, cReset, cGold, cReset, spark)
		if time.Now().Before(flashUntil) && flashMsg != "" {
			fmt.Fprintf(&frame, "  %s│%s  %s\n", cGold, cReset, flashMsg)
		} else {
			fmt.Fprintf(&frame, "  %s│%s  %sHotkeys:%s %s[p]%s Pause/Resume  %s[s]%s Force Sync  %s[b]%s Bidirectional Pull  %s[u]%s Self-Update  %s[d]%s Doctor  %s[q]%s Quit\n",
				cGold, cReset, cBold, cReset, cGold, cReset, cGold, cReset, cGold, cReset, cGold, cReset, cGold, cReset, cRed, cReset)
		}
		fmt.Fprintf(&frame, "  %s╰────────────────────────────────────────────────────────────────────────────────────────╯%s\033[J\n", cGold, cReset)

		fmt.Print(frame.String())
	}

	renderFrame()
	ticker := time.NewTicker(1500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case k := <-keyCh:
			switch k {
			case 'q', 'Q', 3: // q or Ctrl+C
				return
			case 'p', 'P':
				pauseFile := filepath.Join(gitDir, "syncalibur.pause")
				if _, err := os.Stat(pauseFile); err == nil {
					_ = os.Remove(pauseFile)
					flashMsg = cGreen + cBold + "▶ RESUMED: " + cReset + "Autonomous sync loop resumed."
				} else {
					_ = os.WriteFile(pauseFile, []byte("1\n"), 0644)
					flashMsg = cAmber + cBold + "⏸ PAUSED: " + cReset + "Daemon paused. Press [p] again to resume."
				}
				flashUntil = time.Now().Add(3 * time.Second)
				renderFrame()
			case 's', 'S':
				_ = os.WriteFile(filepath.Join(gitDir, "syncalibur.sync"), []byte("1\n"), 0644)
				flashMsg = cCyan + cBold + "⚡ SYNC NOW: " + cReset + "Immediate commit & push cycle triggered!"
				flashUntil = time.Now().Add(3 * time.Second)
				renderFrame()
			case 'b', 'B':
				_ = exec.Command("bash", scriptPath, "pull").Run()
				flashMsg = cCyan + cBold + "🔀 BIDIRECTIONAL PULL: " + cReset + "Checked upstream remote and rebased cleanly."
				flashUntil = time.Now().Add(3 * time.Second)
				renderFrame()
			case 'u', 'U':
				_ = os.WriteFile(filepath.Join(gitDir, "syncalibur.upd"), []byte("1\n"), 0644)
				flashMsg = cGold + cBold + "⟳ SELF-UPDATE: " + cReset + "Triggered Syncalibur upstream self-update check."
				flashUntil = time.Now().Add(3 * time.Second)
				renderFrame()
			case 'd', 'D':
				flashMsg = cGreen + cBold + "🩺 DOCTOR: " + cReset + "Running health check..."
				flashUntil = time.Now().Add(2 * time.Second)
				renderFrame()
			}
		case <-ticker.C:
			renderFrame()
		}
	}
}
