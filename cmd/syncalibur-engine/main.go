// ==============================================================================
// Syncalibur Go Engine (cmd/syncalibur-engine/main.go)
// ------------------------------------------------------------------------------
// High-performance, zero-dependency Go engine powering Syncalibur's:
//   1. Native 2x2 Subpixel TrueColor PNG-to-ANSI Quadrant Renderer (`render-logo`)
//   2. Semantic Conventional Commit Classifier + Diff-Stat + AI Hook (`commit-msg`)
//   3. Multi-Repository Fleet Scanner & Parallel Status Probe (`fleet-scan`, `fleet-status`)
//   4. 7-Day Commit Velocity Sparkline Generator (`sparkline`)
//   5. Multi-Channel Webhook Dispatcher for Discord, Slack & Telegram (`webhook`)
//   6. Interactive Live TUI Control Center (`ui`)
// ==============================================================================
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: syncalibur-engine <render-logo|commit-msg|sparkline|fleet-scan|fleet-status|webhook|ui> [args...]")
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
		if len(os.Args) < 7 {
			os.Exit(1)
		}
		_ = sendWebhook(os.Args[2], os.Args[3], os.Args[4], os.Args[5], os.Args[6], os.Args[7], os.Args[8])
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
// 1. Native 2x2 Subpixel TrueColor PNG-to-ANSI Quadrant Renderer
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

	type rgbaPx struct{ r, g, b, a int }
	grid := make([][]rgbaPx, sh)
	for sy := 0; sy < sh; sy++ {
		grid[sy] = make([]rgbaPx, sw)
		for sx := 0; sx < sw; sx++ {
			// Area-averaging super-sampling + contrast & saturation boost
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
			// Contrast (1.22) & saturation (1.25) enhancement
			lum := 0.299*r + 0.587*g + 0.114*bl
			r = lum + (r-lum)*1.25
			g = lum + (g-lum)*1.25
			bl = lum + (bl-lum)*1.25
			r = (r-128.0)*1.22 + 128.0
			g = (g-128.0)*1.22 + 128.0
			bl = (bl-128.0)*1.22 + 128.0
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
					c1 := avgColor(g1)
					c2 := avgColor(g2)
					errVal := colorSqErr(g1, c1) + colorSqErr(g2, c2)
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

type rgbTuple [3]int

func avgColor[T interface{ r, g, b, a int }](list []T) [3]int {
	return [3]int{0, 0, 0}
}
